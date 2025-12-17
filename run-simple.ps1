Write-Host "=== Ghost Browser Launcher ===" -ForegroundColor Green
Write-Host ""

if (Test-Path "ghost-browser-api-v3.exe") {
    Write-Host "Starting Ghost Browser API v3..." -ForegroundColor Yellow
    Write-Host "Access: http://localhost:8080" -ForegroundColor Cyan
    Write-Host ""
    .\ghost-browser-api-v3.exe
} elseif (Test-Path "ghost-browser-v4/main.go") {
    Write-Host "Starting Ghost Browser v4..." -ForegroundColor Yellow
    Set-Location ghost-browser-v4
    go run main.go
} else {
    Write-Host "No executable found!" -ForegroundColor Red
    Write-Host "Try: .\ghost-browser-api-v3.exe" -ForegroundColor Cyan
}