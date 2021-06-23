package sdk

import (
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const ErrorInvalidLeaseTime = "Invalid lease time"
const ErrorKeyIsRequired string = "Key is required"
const ErrorValueIsRequired string = "Value is required"
const ErrorKeyFormat string = "Key format error"

func sanitizeString(s string) string {
	re := regexp.MustCompile(`[&<>]`)
	z := re.ReplaceAllStringFunc(s, sanitizeStringToken)
	return z
}

func sanitizeStringToken(token string) string {
	return fmt.Sprintf("\\u00%x", int([]rune(token)[0]))
}

func encodeSafe(s string) string {
	return url.QueryEscape(s)
}

func validateKey(key string) error {
	if strings.Contains(key, "/") {
		return fmt.Errorf(ErrorKeyFormat)
	}
	return nil
}

func makeRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
