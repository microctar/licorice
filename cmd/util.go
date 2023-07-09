package cmd

import (
	"encoding/base64"
	"unsafe"
)

func RawURLEncoding(url string) string {
	return base64.RawURLEncoding.EncodeToString(unsafe.Slice(unsafe.StringData(url), len(url)))
}
