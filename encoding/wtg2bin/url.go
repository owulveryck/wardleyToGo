package wtg2bin

import (
	"encoding/base64"
	"fmt"

	"github.com/owulveryck/wardleyToGo/parser/wtg2"
)

var b64 = base64.URLEncoding.WithPadding(base64.NoPadding)

// EncodeURL serializes a Document to a URL-safe string (binary + deflate + base64url, no padding).
func EncodeURL(doc *wtg2.Document) (string, error) {
	data, err := Encode(doc)
	if err != nil {
		return "", err
	}
	return b64.EncodeToString(data), nil
}

// DecodeURL deserializes a URL-safe string back into a Document.
func DecodeURL(s string) (*wtg2.Document, error) {
	data, err := b64.DecodeString(s)
	if err != nil {
		data, err = base64.URLEncoding.DecodeString(s)
		if err != nil {
			return nil, fmt.Errorf("base64url decode: %w", err)
		}
	}
	return Decode(data)
}
