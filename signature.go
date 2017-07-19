package krakenapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"

	"net/url"
)

func getSha256(input []byte) []byte {
	sha := sha256.New()
	sha.Write(input)

	return sha.Sum(nil)
}

func getHMacSha512(message, secret []byte) []byte {
	mac := hmac.New(sha512.New, secret)
	mac.Write(message)

	return mac.Sum(nil)
}

func createKrakenSignature(url_path string, values url.Values, secret []byte) string {
	shasum := getSha256([]byte(values.Get("nonce") + values.Encode()))
	macsum := getHMacSha512(append([]byte(url_path), shasum...), secret)

	return base64.StdEncoding.EncodeToString(macsum)
}
