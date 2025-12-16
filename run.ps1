#!/usr/bin/env pwsh

Write-Host "=== Ghost Browser Launcher ===" -ForegroundColor Green

# Check if executable exists
if (Test-Path "ghost-browser-backend.exe") {
    Write-Host "Starting Ghost Browser backend..." -ForegroundColor Yellow
    Write-Host "Note: This is backend-only mode. For full UI, install Node.js and run 'wails dev'" -ForegroundColor Cyan
    Write-Host "Press Ctrl+C to stop the application" -ForegroundColor Gray
    Write-Host ""
    
    # Run the backend
    .\ghost-browser-backend.exe
} elseif (Test-Path "build/bin/Ghost Browser.exe") {
    Write-Host "Starting Ghost Browser (full version)..." -ForegroundColor Yellow
    & "build/bin/Ghost Browser.exe"
} else {
    Write-Host "Ghost Browser executable not found!" -ForegroundColor Red
    Write-Host "Please run the build script first:" -ForegroundColor Yellow
    Write-Host "  powershell -ExecutionPolicy Bypass -File build_and_test.ps1" -ForegroundColor White
    Write-Host ""
    Write-Host "Or build manually:" -ForegroundColor Yellow
    Write-Host "  go build -o ghost-browser-backend.exe ./cmd/ghost" -ForegroundColor White
    exit 1
}