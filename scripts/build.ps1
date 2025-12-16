# Ghost Browser Build Script for Windows
# Run with: powershell -ExecutionPolicy Bypass -File scripts/build.ps1

Write-Host "=== Ghost Browser Build ===" -ForegroundColor Cyan

# Check requirements
Write-Host "`nChecking requirements..." -ForegroundColor Yellow

# Check Go
if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: Go is not installed" -ForegroundColor Red
    exit 1
}
Write-Host "Go: $(go version)" -ForegroundColor Green

# Check Node
if (!(Get-Command node -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: Node.js is not installed" -ForegroundColor Red
    exit 1
}
Write-Host "Node: $(node --version)" -ForegroundColor Green

# Check Wails
if (!(Get-Command wails -ErrorAction SilentlyContinue)) {
    Write-Host "Wails not found. Installing..." -ForegroundColor Yellow
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
}
Write-Host "Wails: $(wails version)" -ForegroundColor Green

# Install Go dependencies
Write-Host "`nInstalling Go dependencies..." -ForegroundColor Yellow
go mod tidy

# Install frontend dependencies
Write-Host "`nInstalling frontend dependencies..." -ForegroundColor Yellow
Push-Location frontend
npm install
Pop-Location

# Build
Write-Host "`nBuilding application..." -ForegroundColor Yellow
wails build

if ($LASTEXITCODE -eq 0) {
    Write-Host "`n=== Build Complete ===" -ForegroundColor Green
    Write-Host "Executable: build/bin/GhostBrowser.exe" -ForegroundColor Cyan
} else {
    Write-Host "`n=== Build Failed ===" -ForegroundColor Red
    exit 1
}
