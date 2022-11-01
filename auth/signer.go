package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"strings"

	"../../../antcloud-golang-sdk/sdk/utils"
)

// 签名
func Sign(stringToSign, secret string) string {
	return ShaHmac1(stringToSign, secret)
}

// 签名
func SignParams(params map[string]string, secret string) string {
	stringToSign := utils.GetUrlFormedMap(params)
	stringToSign = strings.Replace(stringToSign, "+", "%20", -1)
	stringToSign = strings.Replace(stringToSign, "*", "%2A", -1)
	stringToSign = strings.Replace(stringToSign, "%7E", "~", -1)
	return Sign(stringToSign, secret)
}

func ShaHmac1(source, secret string) string {
	key := []byte(secret)
	hmacCalculator := hmac.New(sha1.New, key)
	hmacCalculator.Write([]byte(source))
	signedBytes := hmacCalculator.Sum(nil)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)
	return signedString
}
