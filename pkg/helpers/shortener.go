package helpers

import gonanoid "github.com/matoous/go-nanoid"

var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func GenToken(length int) (string, error) {
	token, err := gonanoid.Generate(alphabet, length)
	return token, err
}
