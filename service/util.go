package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"time"
)

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

// generateSignedPostPolicyV4 generates a signed post policy.
func generateSignedPostPolicyV4(bucket, objectName string, redirectURL string) (*storage.PostPolicyV4, error) {
	metadata := map[string]string{
		"x-goog-meta-test": "data",
	}

	opts := &storage.PostPolicyV4Options{
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(30 * time.Minute),
		Fields: &storage.PolicyV4Fields{
			Metadata:               metadata,
			RedirectToURLOnSuccess: redirectURL + "/redirect",
		},
	}

	policy, err := storage.GenerateSignedPostPolicyV4(bucket, objectName, opts)
	if err != nil {
		return nil, fmt.Errorf("storage.GenerateSignedPostPolicyV4: %v", err)
	}

	return policy, nil
}

func generateV4GetObjectSignedURL(bucket, object string, conf *jwt.Config, ) (string, error) {

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}
	u, err := storage.SignedURL(bucket, object, opts)
	if err != nil {
		return "", fmt.Errorf("storage.SignedURL: %v", err)
	}

	return u, nil
}
