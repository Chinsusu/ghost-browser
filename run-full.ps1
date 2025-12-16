#!/usr/bin/env pwsh

Write-Host "=== Ghost Browser Full Application ===" -ForegroundColor Green

# Check if full executable exists
if (Test-Path "ghost-browser-full.exe") {
    Write-Host "Starting Ghost Browser (Full Version with UI)..." -ForegroundColor Yellow
    Write-Host "This includes both backend and frontend." -ForegroundColor Cyan
    Write-Host "Press Ctrl+C to stop the application" -ForegroundColor Gray
    Write-Host ""
    
    # Run the full application
    .\ghost-browser-full.exe
} else {
    Write-Host "Full Ghost Browser executable not found!" -ForegroundColor Red
    Write-Host "Please build the full version first:" -ForegroundColor Yellow
    Write-Host "  1. Make sure Node.js is installed" -ForegroundColor White
    Write-Host "  2. cd frontend && npm install && npm run build && cd .." -ForegroundColor White
    Write-Host "  3. go build -o ghost-browser-full.exe ./cmd/ghost" -ForegroundColor White
    Write-Host ""
    Write-Host "Or run the automated build script:" -ForegroundColor Yellow
    Write-Host "  powershell -ExecutionPolicy Bypass -File build_and_test.ps1" -ForegroundColor White
    exit 1
}