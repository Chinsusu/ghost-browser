#!/usr/bin/env pwsh

Write-Host "=== Ghost Browser Desktop Build Script ===" -ForegroundColor Green

Write-Host "`n1. Building frontend..." -ForegroundColor Yellow
Set-Location frontend
npm run build
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Frontend build failed" -ForegroundColor Red
    exit 1
}
Set-Location ..
Write-Host "✅ Frontend built successfully" -ForegroundColor Green

Write-Host "`n2. Building desktop app with proper tags..." -ForegroundColor Yellow
go build -tags desktop,production -ldflags "-w -s" -o ghost-browser-desktop.exe .
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Desktop build failed" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Desktop app built successfully" -ForegroundColor Green

Write-Host "`n🎉 Build completed!" -ForegroundColor Green
Write-Host "Run: .\ghost-browser-desktop.exe" -ForegroundColor Cyan