package helpers

import gonanoid "github.com/matoous/go-nanoid"

var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func GenToken() (string, error) {
	token, err := gonanoid.Generate(alphabet, 8)
	return token, err
}
