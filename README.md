# tosconv

A fast CLI tool for converting TouchOSC `.tosc` files to/from XML.

## What are .tosc files?

TouchOSC layouts are stored as `.tosc` files, which are simply **zlib-compressed XML**. This tool lets you:

- Extract the XML to edit layouts manually or with scripts
- Recompress modified XML back to `.tosc` format

## Installation

### From source

```bash
go install github.com/conormkelly/tosc-converter/cmd/tosconv@latest
```

### From releases

Download the appropriate binary from the [releases page](https://github.com/conormkelly/tosc-converter/releases).

## Usage

```bash
# Decompress .tosc to XML
tosconv layout.tosc                # → layout.xml
tosconv -p layout.tosc             # → layout.xml (pretty-printed)
tosconv -o out.xml layout.tosc     # → out.xml

# Compress XML to .tosc
tosconv layout.xml                 # → layout.tosc
tosconv -o custom.tosc layout.xml

# Force overwrite existing file
tosconv -f layout.tosc
```

The operation is automatically determined by the input file extension:

- `.tosc` → decompress to XML
- `.xml` → compress to .tosc

### Flags

| Flag | Description |
|------|-------------|
| `-o, --output` | Output file path (default: derived from input) |
| `-f, --force` | Overwrite existing output file |
| `-p, --pretty` | Pretty-print XML output (decompress only) |
| `-v, --version` | Show version |
| `-h, --help` | Show help |

## Building

```bash
# Build for current platform
make build

# Install to GOPATH/bin
make install

# Run tests
make test

# Cross-compile for all platforms
make release
```

## License

MIT
