package validation

import (
	"encoding/base64"
	"os"
)

func GetTestPrivateKey() string {
	privateKey, err := base64.StdEncoding.DecodeString(os.Getenv("TEST_PRIVATE_KEY_BASE64"))
	if err != nil {
		panic(err)
	}
	return string(privateKey)
}

func GetTestPublicKey() string {
	pubKey, err := base64.StdEncoding.DecodeString(os.Getenv("TEST_PUBLIC_KEY_BASE64"))
	if err != nil {
		panic(err)
	}
	return string(pubKey)
}
