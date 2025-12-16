#!/usr/bin/env pwsh

Write-Host "=== Ghost Browser Launcher ===" -ForegroundColor Green
Write-Host ""

# Check available executables
$backendOnly = Test-Path "ghost-browser-backend.exe"
$apiServer = Test-Path "ghost-browser-api.exe"
$wailsApp = Test-Path "ghost-browser-wails.exe"

Write-Host "Available versions:" -ForegroundColor Yellow
if ($backendOnly) { Write-Host "✅ Backend Only (ghost-browser-backend.exe)" -ForegroundColor Green }
if ($apiServer) { Write-Host "✅ API Server + Web UI (ghost-browser-api.exe)" -ForegroundColor Green }
if ($wailsApp) { Write-Host "⚠️  Wails Desktop App (ghost-browser-wails.exe) - May have issues" -ForegroundColor Yellow }

Write-Host ""

if ($apiServer) {
    Write-Host "🚀 Starting Ghost Browser API Server..." -ForegroundColor Cyan
    Write-Host "This will start both the backend API and serve the web UI" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Once started, you can access:" -ForegroundColor Yellow
    Write-Host "  📱 Web UI: http://localhost:8080" -ForegroundColor White
    Write-Host "  🔧 API Health: http://localhost:8080/api/health" -ForegroundColor White
    Write-Host "  👤 Profiles API: http://localhost:8080/api/profiles" -ForegroundColor White
    Write-Host "  🌐 Proxies API: http://localhost:8080/api/proxies" -ForegroundColor White
    Write-Host ""
    Write-Host "Press Ctrl+C to stop the server" -ForegroundColor Gray
    Write-Host ""
    
    .\ghost-browser-api.exe
} elseif ($backendOnly) {
    Write-Host "🔧 Starting Ghost Browser Backend Only..." -ForegroundColor Cyan
    Write-Host "This is the backend without UI for testing purposes" -ForegroundColor Gray
    Write-Host ""
    
    .\ghost-browser-backend.exe
} else {
    Write-Host "❌ No Ghost Browser executable found!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please build the application first:" -ForegroundColor Yellow
    Write-Host "  1. For API Server: go build -o ghost-browser-api.exe ./cmd/api" -ForegroundColor White
    Write-Host "  2. For Backend Only: go build -o ghost-browser-backend.exe ./cmd/ghost" -ForegroundColor White
    Write-Host ""
    Write-Host "Or run the automated build script:" -ForegroundColor Yellow
    Write-Host "  powershell -ExecutionPolicy Bypass -File build_and_test.ps1" -ForegroundColor White
    exit 1
}