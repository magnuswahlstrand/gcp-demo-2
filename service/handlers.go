package main

import (
	"bytes"
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/yeqown/go-qrcode"
	"log"
	"net/http"
	"strings"
)

func qrCodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("qr_code_handler")

	oid, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("uuid.Parse: %v", err)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	log.Println("yes")

	f, ok := files[oid]
	if !ok {
		log.Printf("missing file: %v", oid)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	log.Println("yes2")
	conf, err := getConf()
	if err != nil {
		log.Printf("getConf: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("yes3")
	url, err := generateV4GetObjectSignedURL(bucketName, f.objectName, conf)
	if err != nil {
		log.Printf("generateV4GetObjectSignedURL: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("yes4")
	qrc, err := qrcode.New(url, qrcode.WithQRWidth(5))
	if err != nil {
		log.Printf("could not generate QRCode: %v", err)
		return
	}
	log.Println("yes5", url)

	var b bytes.Buffer

	log.Println("yes6", b.Len(), b.Cap())
	b64writer := base64.NewEncoder(base64.StdEncoding, &b)
	log.Println("yes7", b.Len(), b.Cap())
	if err := qrc.SaveTo(b64writer); err != nil {
		log.Printf("qrc.SaveTo: %v\n", err)
		return
	}
	log.Println("yes8", b.Len(), b.Cap())

	data := map[string]string{
		"Base64EncodedImage": b.String(),
		//"Base64EncodedImage": "yeah",
	}
	b64writer.Close()
	log.Println("yes10")

	log.Println("yes7")
	log.Println("yes7")
	if err := templateQRCode.Execute(w, data); err != nil {
		log.Printf("templateQRCode.Execute: %v", err)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	log.Println("yes8")

	log.Println()
}

func helloFormHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("form_handler")
	log.Println("host", r.Host, "host_url", r.URL.Host)

	// TODO: Find another way of determining the service URL. This is not very robust or secure
	var redirectURL string
	if strings.HasPrefix(r.Host, "localhost") {
		redirectURL = "http://" + r.Host
	} else {
		redirectURL = "https://" + r.Host
	}

	conf, err := getConf()
	if err != nil {
		log.Printf("getConf: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	oid := uuid.New()
	policy, err := generateSignedPostPolicyV4(w, bucketName, conf, redirectURL, oid)
	if err != nil {
		log.Printf("generateSignedPostPolicyV4: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files[oid] = file{
		oid: oid,
	}
	log.Println("saving file id", oid.String())
	log.Println(policy.URL)
	log.Println()
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("redirect_handler")

	bucket := r.URL.Query().Get("bucket")
	if bucket != bucketName {
		log.Printf("unexpected bucket name: %s", bucket)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	key := r.URL.Query().Get("key")
	if key == "" {
		log.Printf("key missing")
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	k := strings.Split(key, "/")
	if len(k) != 2 {
		log.Printf("invalid length key: %d", len(k))
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	oid, err := uuid.Parse(k[0])
	if err != nil {
		log.Printf("uuid.Parse: %v", err)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	f, ok := files[oid]
	if !ok {
		log.Printf("missing file: %v", oid)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	f.objectName = key
	files[oid] = f

	log.Println("file found!", f.oid, ",", f.objectName)
	//
	//b := strings.Builder{}
	//b.WriteString(fmt.Sprintf("url: %v\n", r.URL))
	//b.WriteString(fmt.Sprintf("raw_query: %v\n", r.URL.RawQuery))
	//b.WriteString(fmt.Sprintf("query: %v\n", r.URL.Query()))
	//all := b.String()
	//
	//var bd string
	//b2, err := ioutil.ReadAll(r.Body)
	//if err == nil {
	//	bd = string(b2)
	//}
	//
	//log.Println("body", bd)
	log.Println("redirect with url")
	//w.Write([]byte("qr_code?id=" + oid.String()))
	http.Redirect(w, r, "../qr_code?id="+oid.String(), http.StatusFound)
	//log.Println()
}
