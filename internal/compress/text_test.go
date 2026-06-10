package compress

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestBase85RoundTrip(t *testing.T) {
	input := `
component CDN : IV.5 (buy)
component Compute : III.2 (outsource)
anchor User : IV.5
component Engine : II.7 !! >> III.5
CDN -> Compute
User -> CDN
`
	doc := parseWTG2(t, input)

	encoded, err := CompressBase85(doc)
	if err != nil {
		t.Fatalf("CompressBase85: %v", err)
	}
	t.Logf("Base85 output (%d chars): %s", len(encoded), encoded)

	got, err := DecompressBase85(encoded)
	if err != nil {
		t.Fatalf("DecompressBase85: %v", err)
	}
	if !reflect.DeepEqual(doc, got) {
		t.Errorf("round-trip mismatch")
	}
}

func TestBase64URLRoundTrip(t *testing.T) {
	input := `
component CDN : IV.5 (buy)
component Compute : III.2 (outsource)
anchor User : IV.5
component Engine : II.7 !! >> III.5
CDN -> Compute
User -> CDN
`
	doc := parseWTG2(t, input)

	encoded, err := CompressBase64URL(doc)
	if err != nil {
		t.Fatalf("CompressBase64URL: %v", err)
	}
	t.Logf("Base64URL output (%d chars): %s", len(encoded), encoded)

	got, err := DecompressBase64URL(encoded)
	if err != nil {
		t.Fatalf("DecompressBase64URL: %v", err)
	}
	if !reflect.DeepEqual(doc, got) {
		t.Errorf("round-trip mismatch")
	}
}

func TestBase85CharacterSet(t *testing.T) {
	doc := parseWTG2(t, "title: Test\ncomponent A : II.5\n")
	encoded, err := CompressBase85(doc)
	if err != nil {
		t.Fatalf("CompressBase85: %v", err)
	}
	for i, c := range encoded {
		if !strings.ContainsRune(z85Alphabet, c) {
			t.Errorf("invalid char %q at position %d", c, i)
		}
	}
}

func TestBase64URLNoEquals(t *testing.T) {
	doc := parseWTG2(t, "title: Test\ncomponent A : II.5\n")
	encoded, err := CompressBase64URL(doc)
	if err != nil {
		t.Fatalf("CompressBase64URL: %v", err)
	}
	if strings.Contains(encoded, "=") {
		t.Errorf("base64url output should have no padding, got: %s", encoded)
	}
}

func TestDecompressBase85Invalid(t *testing.T) {
	_, err := DecompressBase85("not valid compressed data at all!!!")
	if err == nil {
		t.Error("expected error for invalid input")
	}
}

func TestDecompressBase64URLInvalid(t *testing.T) {
	_, err := DecompressBase64URL("not-valid-data!!!")
	if err == nil {
		t.Error("expected error for invalid input")
	}
}

func TestTextEncodingExampleFiles(t *testing.T) {
	files := []string{
		"../../wtg2/example.wtg2",
		"../../wtg2/full_example.wtg2",
	}

	for _, path := range files {
		t.Run(path, func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Skipf("cannot read %s: %v", path, err)
			}

			doc := parseWTG2(t, string(data))
			originalSize := len(data)

			b85, err := CompressBase85(doc)
			if err != nil {
				t.Fatalf("CompressBase85: %v", err)
			}
			b64, err := CompressBase64URL(doc)
			if err != nil {
				t.Fatalf("CompressBase64URL: %v", err)
			}

			binData, err := CompressBytes(doc)
			if err != nil {
				t.Fatalf("CompressBytes: %v", err)
			}

			t.Logf("%-30s  original: %5d bytes", path, originalSize)
			t.Logf("  binary:    %5d bytes  (%4.1f%%)", len(binData), 100*float64(len(binData))/float64(originalSize))
			t.Logf("  base85:    %5d chars  (%4.1f%%)", len(b85), 100*float64(len(b85))/float64(originalSize))
			t.Logf("  base64url: %5d chars  (%4.1f%%)", len(b64), 100*float64(len(b64))/float64(originalSize))

			got85, err := DecompressBase85(b85)
			if err != nil {
				t.Fatalf("DecompressBase85: %v", err)
			}
			if !reflect.DeepEqual(doc, got85) {
				t.Errorf("base85 round-trip mismatch")
			}

			got64, err := DecompressBase64URL(b64)
			if err != nil {
				t.Fatalf("DecompressBase64URL: %v", err)
			}
			if !reflect.DeepEqual(doc, got64) {
				t.Errorf("base64url round-trip mismatch")
			}
		})
	}
}
