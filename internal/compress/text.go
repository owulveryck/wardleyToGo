package compress

import (
	"encoding/base64"
	"fmt"

	"github.com/owulveryck/wardleyToGo/parser/wtg2"
)

var b64url = base64.URLEncoding.WithPadding(base64.NoPadding)

// CompressBase85 compresses a WTG2 Document and encodes it as a Base85 (Z85) string.
// Suitable for copy-paste, embedding in documents, or QR codes. Not URL-safe.
func CompressBase85(doc *wtg2.Document) (string, error) {
	data, err := CompressBytes(doc)
	if err != nil {
		return "", fmt.Errorf("compress: %w", err)
	}
	return Base85Encode(data), nil
}

// DecompressBase85 decodes a Base85 (Z85) string and decompresses it into a WTG2 Document.
func DecompressBase85(s string) (*wtg2.Document, error) {
	data, err := Base85Decode(s)
	if err != nil {
		return nil, fmt.Errorf("base85 decode: %w", err)
	}
	return DecompressBytes(data)
}

// CompressBase64URL compresses a WTG2 Document and encodes it as a URL-safe
// base64 string (no padding). Suitable for URL fragments and query parameters.
func CompressBase64URL(doc *wtg2.Document) (string, error) {
	data, err := CompressBytes(doc)
	if err != nil {
		return "", fmt.Errorf("compress: %w", err)
	}
	return b64url.EncodeToString(data), nil
}

// DecompressBase64URL decodes a base64url string and decompresses it into a WTG2 Document.
func DecompressBase64URL(s string) (*wtg2.Document, error) {
	data, err := b64url.DecodeString(s)
	if err != nil {
		data, err = base64.URLEncoding.DecodeString(s)
		if err != nil {
			return nil, fmt.Errorf("base64url decode: %w", err)
		}
	}
	return DecompressBytes(data)
}
