package utils

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)


type result struct {
	hash []byte
	err error
}
func HashPassword(password string) (string, error) {
	resultChan := make(chan result,1)

	go func(){
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		resultChan <- result{hash: hash, err: err}
	}()

	select {
	case result := <- resultChan:
		return string(result.hash), result.err
	case <- time.After(5 * time.Second):
		return "", fmt.Errorf("password hashing timed out")
	}
	
}


