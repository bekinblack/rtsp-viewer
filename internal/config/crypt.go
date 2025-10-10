package config

import "encoding/base64"

func Encode(plainText string) string {
	return base64.StdEncoding.EncodeToString([]byte(plainText))
}

func Decode(encodedText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
