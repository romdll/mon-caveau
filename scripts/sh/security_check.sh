#!/bin/bash

OUTPUT_FILE="security_report.txt"

echo "Installing necessary Go tools..."
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

echo "Running security scans..." > "$OUTPUT_FILE"

echo "Running gosec..." | tee -a "$OUTPUT_FILE"
gosec ./... >> "$OUTPUT_FILE"

echo "Running govulncheck..." | tee -a "$OUTPUT_FILE"
govulncheck -show verbose ./... >> "$OUTPUT_FILE"

echo "Running staticcheck..." | tee -a "$OUTPUT_FILE"
staticcheck ./... >> "$OUTPUT_FILE"

echo "Security scans completed. Check $OUTPUT_FILE for details."