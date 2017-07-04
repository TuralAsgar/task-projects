package main

import (
	"crypto/rand"
	"fmt"
)

const (
	randomBytes = 15
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Random(bytes int) (string, error) {
	b, err := Bytes(bytes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
