# Ghost Browser Test v3 - Using ChromeDP (no download needed)
# This version uses your installed Edge directly

Write-Host "=== Ghost Browser Test v3 ===" -ForegroundColor Cyan
Write-Host "Using ChromeDP with installed Edge (no download)" -ForegroundColor Yellow
Write-Host ""

# Get dependencies
Write-Host "Getting dependencies..." -ForegroundColor Gray
go mod tidy

# Run
Write-Host ""
go run main.go
