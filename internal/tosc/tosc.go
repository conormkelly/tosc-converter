// Package tosc provides functions for converting TouchOSC .tosc files to/from XML.
// .tosc files are zlib-compressed XML.
package tosc

import (
	"bytes"
	"compress/zlib"
	"encoding/xml"
	"fmt"
	"io"
)

// Decompress takes zlib-compressed data and returns the decompressed XML.
func Decompress(input []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer r.Close()

	output, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress: %w", err)
	}

	return output, nil
}

// Compress takes XML data and returns zlib-compressed data.
func Compress(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)

	_, err := w.Write(input)
	if err != nil {
		return nil, fmt.Errorf("failed to compress: %w", err)
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize compression: %w", err)
	}

	return buf.Bytes(), nil
}

// PrettyPrint formats XML with indentation.
func PrettyPrint(input []byte) ([]byte, error) {
	// First decode to verify it's valid XML and normalize it
	decoder := xml.NewDecoder(bytes.NewReader(input))

	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to parse XML: %w", err)
		}
		if err := encoder.EncodeToken(token); err != nil {
			return nil, fmt.Errorf("failed to encode XML: %w", err)
		}
	}

	if err := encoder.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush XML: %w", err)
	}

	// Add trailing newline
	result := buf.Bytes()
	if len(result) > 0 && result[len(result)-1] != '\n' {
		result = append(result, '\n')
	}

	return result, nil
}
