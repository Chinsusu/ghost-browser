#!/usr/bin/env pwsh

Write-Host "🧪 === Ghost Browser Launcher Test ===" -ForegroundColor Green
Write-Host ""

Write-Host "🔧 Building test executable..." -ForegroundColor Yellow
go build -o test-browser-launcher.exe test-browser-launcher.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Build failed" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Test executable built successfully" -ForegroundColor Green

Write-Host "`n🚀 Running browser launcher test..." -ForegroundColor Yellow
Write-Host "This will:" -ForegroundColor Gray
Write-Host "  1. Create a test profile with random fingerprint" -ForegroundColor Gray
Write-Host "  2. Launch Edge browser with spoofed fingerprint" -ForegroundColor Gray
Write-Host "  3. Provide test URLs and instructions" -ForegroundColor Gray
Write-Host "  4. Clean up test data" -ForegroundColor Gray
Write-Host ""

.\test-browser-launcher.exe

Write-Host "`n🧹 Cleaning up..." -ForegroundColor Yellow
Remove-Item -Path "test-browser-launcher.exe" -ErrorAction SilentlyContinue
Remove-Item -Path "test-browser-launcher.go" -ErrorAction SilentlyContinue

Write-Host "Test completed!" -ForegroundColor Green