package main

import (
	"cloud.google.com/go/storage"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/oauth2/jwt"
	"io"
	"time"
)

// generateSignedPostPolicyV4 generates a signed post policy.
func generateSignedPostPolicyV4(w io.Writer, bucket string, conf *jwt.Config, redirectURL string, objectID uuid.UUID) (*storage.PostPolicyV4, error) {
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

	policy, err := storage.GenerateSignedPostPolicyV4(bucket, objectID.String()+"/${filename}", opts)
	if err != nil {
		return nil, fmt.Errorf("storage.GenerateSignedPostPolicyV4: %v", err)
	}

	// Generate the form, using the data from the policy.
	if err = templateUpload.Execute(w, policy); err != nil {
		return policy, fmt.Errorf("executing template: %v", err)
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
