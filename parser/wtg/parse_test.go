package wtg

import (
	"bytes"
	"testing"
)

func TestParseComponent(t *testing.T) {
	p := NewParser()
	test := `
	my component à: {
		cool
		evolution: |.........|.......x.......|.........|......|
	}
	my component: {
		cool
		evolution: |x|..............|.........|......|
	}
	component: {
		cool
		evolution: |....+..............|.........|......|
	}
	`
	p.Parse(bytes.NewBufferString(test))
}
