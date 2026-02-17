package config

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
)

// NewClient creates a Firestore client using the FIREBASE_PROJECT_ID env var.
func NewClient(ctx context.Context) (*firestore.Client, error) {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		// fallback to older env name if present
		projectID = os.Getenv("PROJECT_ID")
	}
	return firestore.NewClient(ctx, projectID)
}
