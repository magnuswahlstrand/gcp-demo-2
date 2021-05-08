package main

import (
	"embed"
	"github.com/google/uuid"
	"golang.org/x/oauth2/jwt"
	"html/template"
	"log"
	"net/http"
	"os"
)

//go:embed templates/*.gohtml
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

	conf *jwt.Config
)

func init() {
	templateUpload = template.Must(template.New("upload.gohtml").ParseFS(templatesFS, "templates/upload.gohtml"))
	templateQRCode = template.Must(template.New("qrcode.gohtml").ParseFS(templatesFS, "templates/qrcode.gohtml"))


	if os.Getenv("BUCKET_NAME") == "" {
		log.Fatal("env var BUCKET_NAME missing")
	}
	bucketName = os.Getenv("BUCKET_NAME")

	if os.Getenv("SERVICE_ACCOUNT_JSON") == "" {
		log.Fatal("env var SERVICE_ACCOUNT_JSON missing")
	}
	secretName = os.Getenv("SERVICE_ACCOUNT_JSON")

	var err error
	conf, err = getConf()
	if err != nil {
		log.Fatal("failed to get credentials", err)
	}
}

func main() {
	http.HandleFunc("/", uploadHandler)
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