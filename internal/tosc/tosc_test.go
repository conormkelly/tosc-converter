package tosc

import (
	"bytes"
	"testing"
)

// Sample XML for testing
var testXML = []byte(`<?xml version='1.0' encoding='UTF-8'?><lexml version="5"><node ID="test" type="GROUP"><properties><property type="s"><key>name</key><value>test</value></property></properties></node></lexml>`)

func TestDecompress(t *testing.T) {
	// First compress some known XML, then decompress it
	compressed, err := Compress(testXML)
	if err != nil {
		t.Fatalf("setup failed - Compress: %v", err)
	}

	// Decompress
	xmlData, err := Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}

	// Should match original
	if !bytes.Equal(xmlData, testXML) {
		t.Errorf("decompressed data doesn't match original")
	}
}

func TestCompress(t *testing.T) {
	// Compress
	compressed, err := Compress(testXML)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	// Should produce non-empty output
	if len(compressed) == 0 {
		t.Error("compressed data is empty")
	}

	// Compressed should be smaller than original (for reasonable XML)
	if len(compressed) >= len(testXML) {
		t.Errorf("compressed size (%d) should be smaller than original (%d)", len(compressed), len(testXML))
	}
}

func TestRoundTrip(t *testing.T) {
	// Compress
	compressed, err := Compress(testXML)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	// Decompress
	decompressed, err := Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}

	// Should match original
	if !bytes.Equal(decompressed, testXML) {
		t.Error("round-trip content doesn't match")
	}

	// Compress again
	recompressed, err := Compress(decompressed)
	if err != nil {
		t.Fatalf("second Compress failed: %v", err)
	}

	// Decompress again
	decompressed2, err := Decompress(recompressed)
	if err != nil {
		t.Fatalf("second Decompress failed: %v", err)
	}

	// Should still match
	if !bytes.Equal(decompressed2, testXML) {
		t.Error("double round-trip content doesn't match")
	}
}

func TestDecompressInvalid(t *testing.T) {
	// Random bytes that aren't zlib
	_, err := Decompress([]byte("not zlib compressed data"))
	if err == nil {
		t.Error("expected error for invalid zlib data")
	}
}

func TestCompressEmpty(t *testing.T) {
	// Empty input should still work
	compressed, err := Compress([]byte{})
	if err != nil {
		t.Fatalf("Compress failed on empty input: %v", err)
	}
	if len(compressed) == 0 {
		t.Error("compressed empty input shouldn't be empty (zlib header)")
	}

	// Should be able to decompress empty
	decompressed, err := Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress empty failed: %v", err)
	}
	if len(decompressed) != 0 {
		t.Error("decompressed empty should be empty")
	}
}

func TestPrettyPrint(t *testing.T) {
	input := []byte(`<?xml version='1.0' encoding='UTF-8'?><root><child>text</child></root>`)

	output, err := PrettyPrint(input)
	if err != nil {
		t.Fatalf("PrettyPrint failed: %v", err)
	}

	// Should contain newlines (formatted)
	if !bytes.Contains(output, []byte("\n")) {
		t.Error("pretty-printed output should contain newlines")
	}

	// Should end with newline
	if len(output) > 0 && output[len(output)-1] != '\n' {
		t.Error("pretty-printed output should end with newline")
	}
}

func TestPrettyPrintInvalid(t *testing.T) {
	// Malformed XML with mismatched tags
	_, err := PrettyPrint([]byte("<root><unclosed>"))
	if err == nil {
		t.Error("expected error for invalid XML")
	}
}
