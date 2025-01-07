@echo off

set OUTPUT_FILE=security_report.txt

echo Installing necessary Go tools...
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

echo Running security scans... > %OUTPUT_FILE%

echo Running gosec...
echo Running gosec... >> %OUTPUT_FILE%
gosec ./... >> %OUTPUT_FILE%

echo Running govulncheck...
echo Running govulncheck... >> %OUTPUT_FILE%
govulncheck ./... >> %OUTPUT_FILE%

echo Running staticcheck...
echo Running staticcheck... >> %OUTPUT_FILE%
staticcheck -show verbose ./... >> %OUTPUT_FILE%

echo Security scans completed. Check %OUTPUT_FILE% for details.