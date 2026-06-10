// wtg2compress compresses and decompresses WTG2 Wardley Map files using
// grammar-aware arithmetic coding.
//
// Usage:
//
//	cat map.wtg2  | wtg2compress                    > map.wtg2c    # compress (binary)
//	cat map.wtg2c | wtg2compress -d                 > map.wtg2     # decompress (binary)
//	cat map.wtg2  | wtg2compress -format base85                    # compress to base85 text
//	cat map.wtg2  | wtg2compress -format base64url                 # compress to base64url text
package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/owulveryck/wardleyToGo/internal/compress"
	"github.com/owulveryck/wardleyToGo/parser/wtg2"
)

func main() {
	decomp := flag.Bool("d", false, "decompress (default: compress)")
	format := flag.String("format", "binary", "output format: binary, base85, base64url")
	flag.Parse()

	if *decomp {
		decompressCmd(*format)
	} else {
		compressCmd(*format)
	}
}

func compressCmd(format string) {
	p, err := wtg2.NewParser(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}

	switch format {
	case "binary":
		if err := compress.Compress(doc, os.Stdout); err != nil {
			log.Fatal(err)
		}
	case "base85":
		s, err := compress.CompressBase85(doc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(s)
	case "base64url":
		s, err := compress.CompressBase64URL(doc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(s)
	default:
		log.Fatalf("unknown format: %s (use binary, base85, or base64url)", format)
	}
}

func decompressCmd(format string) {
	var doc *wtg2.Document
	var err error

	switch format {
	case "binary":
		doc, err = compress.Decompress(bufio.NewReader(os.Stdin))
	case "base85":
		var data []byte
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		doc, err = compress.DecompressBase85(string(data))
	case "base64url":
		var data []byte
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		b64 := base64.URLEncoding.WithPadding(base64.NoPadding)
		var binData []byte
		binData, err = b64.DecodeString(string(data))
		if err != nil {
			log.Fatal(err)
		}
		doc, err = compress.DecompressBytes(binData)
	default:
		log.Fatalf("unknown format: %s (use binary, base85, or base64url)", format)
	}
	if err != nil {
		log.Fatal(err)
	}

	if err := wtg2.Emit(os.Stdout, doc); err != nil {
		log.Fatal(err)
	}
}
