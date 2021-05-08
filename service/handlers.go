package main

import (
	"bytes"
	"cloud.google.com/go/storage"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/yeqown/go-qrcode"
	"net/http"
	"strings"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	policy, err := uploadPolicy(r)
	if err != nil {
		logger.Printf("uploadPolicy: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate the form, using the data from the policy.
	if err = templateUpload.Execute(w, policy); err != nil {
		logger.Printf("templateUpload.Execute: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func uploadPolicy(r *http.Request) (*storage.PostPolicyV4, error) {
	var redirectURL string
	if strings.HasPrefix(r.Host, "localhost") {
		redirectURL = "http://" + r.Host
	} else {
		redirectURL = "https://" + r.Host
	}

	oid := uuid.New()

	// Create signed URL for POST upload
	objectName := fmt.Sprintf("%s/${filename}", oid)
	policy, err := generateSignedPostPolicyV4(bucketName, objectName, redirectURL)
	if err != nil {
		return nil, fmt.Errorf("generateSignedPostPolicyV4: %w", err)
	}

	files[oid] = file{
		oid: oid,
	}
	return policy, nil
}

func qrCodeHandler(w http.ResponseWriter, r *http.Request) {
	qrCodeImage, err := qrCode(r)
	if err != nil {
		logger.Printf("qrCode: %v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	data := map[string]string{"Base64EncodedImage": qrCodeImage}
	if err := templateQRCode.Execute(w, data); err != nil {
		logger.Printf("templateQRCode.Execute: %v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
}


func qrCode(r *http.Request) (string, error){
	oid, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		return "", fmt.Errorf("uuid.Parse: %v", err)
	}

	f, ok := files[oid]
	if !ok {
		return "", fmt.Errorf("missing file: %v", oid)
	}

	url, err := generateV4GetObjectSignedURL(bucketName, f.objectName, conf)
	if err != nil {
		return "", fmt.Errorf("generateV4GetObjectSignedURL: %v", err)
	}

	qrc, err := qrcode.New(url, qrcode.WithQRWidth(5))
	if err != nil {
		return "", fmt.Errorf("could not generate QRCode: %v", err)
	}

	var b bytes.Buffer
	b64writer := base64.NewEncoder(base64.StdEncoding, &b)
	defer b64writer.Close()
	if err := qrc.SaveTo(b64writer); err != nil {
		return "", fmt.Errorf("qrc.SaveTo: %v\n", err)
	}

	return b.String(), nil
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	logger.Println("redirect_handler")
	url, err := redirectURL(r)
	if err != nil {
		logger.Printf("redirectURL: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func redirectURL(r *http.Request) (string, error){
	bucket := r.URL.Query().Get("bucket")
	if bucket != bucketName {
		return "", fmt.Errorf("unexpected bucket name: %s", bucket)
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		return "", fmt.Errorf("key missing")
	}

	k := strings.Split(key, "/")
	if len(k) != 2 {
		return "", fmt.Errorf("invalid length key: %d", len(k))
	}

	oid, err := uuid.Parse(k[0])
	if err != nil {
		return "", fmt.Errorf("uuid.Parse: %v", err)
	}

	f, ok := files[oid]
	if !ok {
		return "", fmt.Errorf("missing file: %v", oid)
	}
	f.objectName = key
	files[oid] = f

	return "../qr_code?id="+oid.String(), nil
}