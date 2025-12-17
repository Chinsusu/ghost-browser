# Ghost Browser - Edge Launcher Test
# Run this script to test if Edge launches correctly with fingerprint spoofing

Write-Host "=== Ghost Browser - Edge Launcher Test ===" -ForegroundColor Cyan
Write-Host ""

# Check Go
if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: Go is not installed" -ForegroundColor Red
    exit 1
}

# Install dependencies
Write-Host "Installing dependencies..." -ForegroundColor Yellow
go mod tidy

# Run test
Write-Host ""
Write-Host "Launching Edge with fingerprint spoofing..." -ForegroundColor Green
Write-Host "Browser will open with 2 tabs:" -ForegroundColor Yellow
Write-Host "  1. browserleaks.com/javascript - Check Navigator values" -ForegroundColor Yellow
Write-Host "  2. CreepJS - Comprehensive fingerprint analysis" -ForegroundColor Yellow
Write-Host ""
Write-Host "Press Ctrl+C to stop" -ForegroundColor Gray
Write-Host ""

go run ./cmd/test-launcher/main.go
