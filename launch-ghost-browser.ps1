#!/usr/bin/env pwsh

Write-Host "🚀 === Ghost Browser - Enterprise Anti-Detection Suite ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "Version: 5.0.0 (Final Release)" -ForegroundColor Green
Write-Host "Status: Production Ready with Advanced ChromeDP Integration" -ForegroundColor Green
Write-Host ""

Write-Host "Available Versions:" -ForegroundColor Yellow
Write-Host "  [1] 🌐 Ghost Browser API Server (Recommended)" -ForegroundColor White
Write-Host "      - Full web interface at http://localhost:8080" -ForegroundColor Gray
Write-Host "      - Profile management, proxy support" -ForegroundColor Gray
Write-Host "      - Advanced ChromeDP browser engine" -ForegroundColor Gray
Write-Host ""
Write-Host "  [2] 🧪 Ghost Browser v4 Standalone" -ForegroundColor White
Write-Host "      - Direct browser launch with random fingerprints" -ForegroundColor Gray
Write-Host "      - Automatic CreepJS testing" -ForegroundColor Gray
Write-Host "      - Advanced spoofing demonstration" -ForegroundColor Gray
Write-Host ""

$choice = Read-Host "Select version (1 or 2)"

switch ($choice) {
    "1" {
        Write-Host ""
        Write-Host "🚀 Starting Ghost Browser API Server..." -ForegroundColor Yellow
        Write-Host "✅ Features: ChromeDP Integration, Advanced Spoofing" -ForegroundColor Green
        Write-Host "🌐 Access: http://localhost:8080" -ForegroundColor Cyan
        Write-Host "Press Ctrl+C to stop" -ForegroundColor Gray
        Write-Host ""
        .\ghost-browser-release.exe
    }
    "2" {
        Write-Host ""
        Write-Host "🧪 Starting Ghost Browser v4 Standalone..." -ForegroundColor Yellow
        Write-Host "✅ Features: Random Fingerprints, CreepJS Testing" -ForegroundColor Green
        Write-Host "🔍 Auto-opens: browserleaks.com + CreepJS" -ForegroundColor Cyan
        Write-Host "Press Ctrl+C to stop" -ForegroundColor Gray
        Write-Host ""
        .\ghost-browser-v4-standalone.exe
    }
    default {
        Write-Host ""
        Write-Host "❌ Invalid choice. Starting API Server (default)..." -ForegroundColor Red
        Write-Host ""
        .\ghost-browser-release.exe
    }
}