package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"embed"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"html/template"
	"log"
	"net/http"
	"os"
)

//go:embed *.gohtml
var templatesFS embed.FS

//go:embed assets/*
var assertsFS embed.FS

type file struct {
	oid        uuid.UUID
	objectName string
}

var (
	files = map[uuid.UUID]file{}

	templateUpload, templateQRCode *template.Template

	bucketName string
	secretName string
)

func init() {
	templateUpload = template.Must(template.New("upload.gohtml").ParseFS(templatesFS, "upload.gohtml"))
	templateQRCode = template.Must(template.New("qrcode.gohtml").ParseFS(templatesFS, "qrcode.gohtml"))


	if os.Getenv("BUCKET_NAME") == "" {
		log.Fatal("env var BUCKET_NAME missing")
	}
	bucketName = os.Getenv("BUCKET_NAME")

	if os.Getenv("SERVICE_ACCOUNT_JSON") == "" {
		log.Fatal("env var SERVICE_ACCOUNT_JSON missing")
	}
	secretName = os.Getenv("SERVICE_ACCOUNT_JSON")
}

func main() {
	http.HandleFunc("/", helloFormHandler)
	http.HandleFunc("/qr_code", qrCodeHandler)
	http.HandleFunc("/redirect", redirectHandler)

	fileServer := http.FileServer(http.FS(assertsFS))
	http.Handle("/assets/", fileServer)

	// PORT environment variable is provided by Cloud Run.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Print("Hello from Cloud Run! The container started successfully and is listening for HTTP requests on $PORT")
	log.Printf("Listening on port %s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func getSecret(name string) (string, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", err
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", err
	}

	return string(result.Payload.Data), nil
}

func getConf() (*jwt.Config, error) {
	secret, err := getSecret(secretName)
	if err != nil {
		return nil, fmt.Errorf("getSecret(%s): %w", secretName, err)

	}

	conf, err := google.JWTConfigFromJSON([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("google.JWTConfigFromJSON: %w", err)
	}
	return conf, nil
}
