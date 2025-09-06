package utils

import (
	"fmt"
	"time"
)

func Retry(attempts int, sleep time.Duration, fn func() error) error {
	for i := 0; i < attempts; i++ {
		if err := fn(); err == nil {
			return nil
		}
	}
	return fmt.Errorf("failed to execute function after %d attempts", attempts)
}