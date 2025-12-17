# Ghost Browser - Fixed Launcher Test v2
# This version injects script BEFORE page load

Write-Host "=== Ghost Browser - Launcher Test v2 ===" -ForegroundColor Cyan
Write-Host "This version uses AddScriptToEvaluateOnNewDocument" -ForegroundColor Yellow
Write-Host ""

go mod tidy 2>$null
go run ./cmd/test-launcher/main.go
