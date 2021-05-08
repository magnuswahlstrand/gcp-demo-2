package main

import (
	"embed"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	logger *log.Logger
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

	logger = log.Default()
}

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", uploadHandler)
	r.Get("/qr_code", qrCodeHandler)
	r.Get("/redirect", redirectHandler)
	r.Handle("/assets/", http.FileServer(http.FS(assertsFS)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Printf("hello from cloud run! the container started successfully and is listening for HTTP requests on %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
