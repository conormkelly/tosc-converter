package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/conormkelly/tosc-converter/internal/tosc"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Define flags
	var (
		outputPath  string
		force       bool
		pretty      bool
		showVersion bool
	)

	flag.StringVar(&outputPath, "o", "", "output file path")
	flag.StringVar(&outputPath, "output", "", "output file path")
	flag.BoolVar(&force, "f", false, "overwrite existing output file")
	flag.BoolVar(&force, "force", false, "overwrite existing output file")
	flag.BoolVar(&pretty, "p", false, "pretty-print XML output (decompress only)")
	flag.BoolVar(&pretty, "pretty", false, "pretty-print XML output (decompress only)")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&showVersion, "version", false, "show version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `tosconv - Convert TouchOSC .tosc files to/from XML

Usage:
  tosconv [flags] <input-file>

The operation is determined by the input file extension:
  .tosc → decompress to XML
  .xml  → compress to .tosc

Flags:
  -o, --output <path>   Output file path (default: derived from input)
  -f, --force           Overwrite existing output file
  -p, --pretty          Pretty-print XML output (decompress only)
  -v, --version         Show version
  -h, --help            Show help

Examples:
  tosconv layout.tosc                # → layout.xml
  tosconv -p layout.tosc             # → layout.xml (pretty-printed)
  tosconv -o out.xml layout.tosc     # → out.xml
  tosconv layout.xml                 # → layout.tosc
  tosconv -f layout.tosc             # Overwrite existing layout.xml
`)
	}

	flag.Parse()

	if showVersion {
		fmt.Printf("tosconv %s\n", version)
		return nil
	}

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return fmt.Errorf("expected exactly one input file")
	}

	inputPath := args[0]

	// Determine operation based on extension
	ext := strings.ToLower(filepath.Ext(inputPath))
	var decompress bool
	switch ext {
	case ".tosc":
		decompress = true
	case ".xml":
		decompress = false
	default:
		return fmt.Errorf("unsupported file extension %q (expected .tosc or .xml)", ext)
	}

	// Derive output path if not specified
	if outputPath == "" {
		base := strings.TrimSuffix(inputPath, ext)
		if decompress {
			outputPath = base + ".xml"
		} else {
			outputPath = base + ".tosc"
		}
	}

	// Check if output file exists
	if !force {
		if _, err := os.Stat(outputPath); err == nil {
			return fmt.Errorf("output file %q already exists (use -f to overwrite)", outputPath)
		}
	}

	// Read input file
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Process
	var output []byte
	if decompress {
		output, err = tosc.Decompress(input)
		if err != nil {
			return fmt.Errorf("failed to decompress: %w", err)
		}
		if pretty {
			output, err = tosc.PrettyPrint(output)
			if err != nil {
				return fmt.Errorf("failed to pretty-print: %w", err)
			}
		}
	} else {
		output, err = tosc.Compress(input)
		if err != nil {
			return fmt.Errorf("failed to compress: %w", err)
		}
	}

	// Write output file
	if err := os.WriteFile(outputPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	// Print output path to stdout
	fmt.Println(outputPath)
	return nil
}
