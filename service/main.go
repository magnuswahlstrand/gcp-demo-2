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

//go:embed index.html
var fs embed.FS

//go:embed assets/*
var assertsFS embed.FS

// templateData provides template parameters.
type templateData struct {
	Service  string
	Revision string
	Secret   string
}

// Variables used to generate the HTML page.
var (
	data  templateData
	data2 map[string]string
	tmpl  *template.Template

	jsonSecret = os.Getenv("SERVICE_ACCOUNT_JSON")
	bucketName = os.Getenv("BUCKET_NAME")

	files = map[uuid.UUID]file{}
)

type file struct {
	oid        uuid.UUID
	objectName string
}

func main() {
	tmpl = template.Must(template.ParseFS(fs, "index.html"))

	data2 = map[string]string{
		"Service": "service",
		//"Secret":     secret,

		"URL":    "url",
		"PutURL": "urlPUT",
	}

	// Define HTTP server.
	//http.HandleFunc("/", helloRunHandler)
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

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", err
	}

	// WARNING: Do not print the secret in a production environment - this snippet
	// is showing how to access the secret material.
	secret := string(result.Payload.Data)
	return secret, nil
}

func getConf() (*jwt.Config, error) {
	secret, err := getSecret(jsonSecret)
	if err != nil {
		return nil, fmt.Errorf("getSecret(%s): %w", jsonSecret, err)

	}

	conf, err := google.JWTConfigFromJSON([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("google.JWTConfigFromJSON: %w", err)
	}
	return conf, nil
}

// form is a template for an HTML form that will use the data from the signed
var form = `
<html>
  <body>
	<form action="{{ .URL }}" method="POST" enctype="multipart/form-data">
			{{- range $name, $value := .Fields }}
			<input name="{{ $name }}" value="{{ $value }}" type="hidden"/>
			{{- end }}
			<input type="file" name="file"/><br />
			<input type="submit" value="Upload File" /><br />
	</form>
  </body>
</html>
`

var qrCode = `
<html>
    <head>
        <title>Testing QR code</title>
    </head>
    <body>
		<div>
		  <img src="data:image/png;base64, {{ .Base64EncodedImage }}" width=400 height=400 />
		</div>
    </body>
</html>
`

// post policy.

var tmpl2 = template.Must(template.New("policyV4").Parse(form))
var tmpl3 = template.Must(template.New("qrCode").Parse(qrCode))
