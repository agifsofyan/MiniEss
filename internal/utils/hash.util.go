package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// MD5 Hash
func SetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Base64StrEncode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64StrDecode(encodedString string) string {
	var decodedByte, _ = base64.StdEncoding.DecodeString(encodedString)
	return string(decodedByte)
}

type Hash struct{}

func Generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func Compare(hash, s string) bool {
	incoming := []byte(s)
	existing := []byte(hash)
	compare := bcrypt.CompareHashAndPassword(existing, incoming)
	return compare == nil
}
