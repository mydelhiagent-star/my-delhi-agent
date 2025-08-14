package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/benbusby/b2"
)

type B2Service struct {
	svc      *b2.Service
	bucketID string
	mu       sync.RWMutex
}

var (
	BackblazeInstance *B2Service
	once              sync.Once
)

type FileResponse struct {
    FileName  string `json:"fileName"`
    UploadURL string `json:"uploadURL"`
    AuthToken string `json:"authToken"`
}

// Initialize starts a background goroutine to initialize B2Service
func InitializeB2Service(ctx context.Context) {
	once.Do(func() {
		BackblazeInstance = &B2Service{}
		go func() {
			for {
				BackblazeInstance.mu.Lock()
				if BackblazeInstance.svc == nil {
					svc, err := newB2Service()
					if err != nil {
						log.Printf("⚠️ Backblaze init failed, retrying in 30s: %v", err)
					} else {
						BackblazeInstance.svc = svc.svc
						BackblazeInstance.bucketID = svc.bucketID
						log.Println("✅ Backblaze initialized")
					}
				}
				BackblazeInstance.mu.Unlock()
				time.Sleep(30 * time.Second)
			}
		}()
	})
}

// newB2Service initializes Backblaze client
func newB2Service() (*B2Service, error) {
	accountID := os.Getenv("B2_ACCOUNT_ID")
	appKey := os.Getenv("B2_APP_KEY")
	bucketID := os.Getenv("B2_BUCKET_ID")

	if accountID == "" || appKey == "" || bucketID == "" {
		return nil, fmt.Errorf("missing B2 credentials or bucket ID")
	}

	svc, auth, err := b2.AuthorizeAccountV2(accountID, appKey)
	if err != nil {
		return nil, err
	}

	if auth.Allowed.BucketID != bucketID {
		return nil, fmt.Errorf("app key not authorized for bucket %s", bucketID)
	}

	return &B2Service{
		svc:      svc,
		bucketID: bucketID,
	}, nil
}

// GenerateSignedUploadURLs safely uses the singleton instance
func GenerateSignedUploadURLs(ctx context.Context, fileNames []string) ([]FileResponse, error) {
	BackblazeInstance.mu.RLock()
	defer BackblazeInstance.mu.RUnlock()

	if BackblazeInstance == nil || BackblazeInstance.svc == nil {
		return nil, fmt.Errorf("Backblaze service not initialized yet")
	}

	var responses []FileResponse
	for _, fileName := range fileNames {
		resp, err := BackblazeInstance.svc.GetUploadURL(BackblazeInstance.bucketID)
		if err != nil {
			return nil, fmt.Errorf("error getting upload URL for %s: %w", fileName, err)
		}
		responses = append(responses, FileResponse{
			FileName:  fileName,
			UploadURL: resp.UploadURL,
			AuthToken: resp.AuthorizationToken,
		})
	}
	return responses, nil
}
