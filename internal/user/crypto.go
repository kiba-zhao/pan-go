package user

import "encoding/base64"

func DecodeSignature(signature string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(signature)
}

func EncodeSignature(signature []byte) string {
	return base64.RawURLEncoding.EncodeToString(signature)
}

func DecodeUserKey(userKey string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(userKey)
}

func EncodeUserKey(userKey []byte) string {
	return base64.RawURLEncoding.EncodeToString(userKey)
}

func DecodePassphrase(passphrase string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(passphrase)
}

func EncodePassphrase(passphrase []byte) string {
	return base64.RawURLEncoding.EncodeToString(passphrase)
}
