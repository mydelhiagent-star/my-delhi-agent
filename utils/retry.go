// utils/retry.go
package utils

import (
    "context"
    "fmt"
    "time"
)

// ← SIMPLE: One function for all retries
func Retry(ctx context.Context, operation func() error) error {
    const maxRetries = 3
    const baseDelay = 100 * time.Millisecond
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := operation()
        
        if err == nil {
            return nil // Success!
        }
        
        // ← RETRY: Only on transient errors
        if  attempt < maxRetries-1 {
            delay := baseDelay * time.Duration(1<<attempt) // Exponential backoff
            time.Sleep(delay)
            continue
        }
        
        return err // Non-retryable error
    }
    
    return fmt.Errorf("operation failed after %d attempts", maxRetries)
}


