#!/usr/bin/env pwsh

Write-Host "🧪 === Direct CDP Browser Test (Pure Go Implementation) ===" -ForegroundColor Green
Write-Host ""

Write-Host "This test completely avoids external libraries by using:" -ForegroundColor Yellow
Write-Host "  • Direct HTTP calls to Chrome DevTools Protocol" -ForegroundColor Gray
Write-Host "  • Native Go HTTP client for CDP communication" -ForegroundColor Gray
Write-Host "  • JavaScript injection via CDP Runtime.evaluate" -ForegroundColor Gray
Write-Host "  • No Rod, no external dependencies" -ForegroundColor Gray
Write-Host ""

Write-Host "🔧 Building CDP direct test executable..." -ForegroundColor Yellow
go build -o test-cdp-direct.exe test-cdp-direct.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Build failed" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Test executable built successfully" -ForegroundColor Green

Write-Host "`n🚀 Running CDP direct test..." -ForegroundColor Yellow
Write-Host "This will:" -ForegroundColor Gray
Write-Host "  1. Create a test profile with random fingerprint" -ForegroundColor Gray
Write-Host "  2. Launch Edge with CDP enabled on port 9222" -ForegroundColor Gray
Write-Host "  3. Connect to CDP via HTTP API" -ForegroundColor Gray
Write-Host "  4. Inject fingerprint spoofing via Runtime.evaluate" -ForegroundColor Gray
Write-Host "  5. Navigate to test pages" -ForegroundColor Gray
Write-Host ""

.\test-cdp-direct.exe

Write-Host "`n🧹 Cleaning up..." -ForegroundColor Yellow
Remove-Item -Path "test-cdp-direct.exe" -ErrorAction SilentlyContinue

Write-Host "CDP direct test completed!" -ForegroundColor Green
Write-Host ""
Write-Host "📊 Please report the results:" -ForegroundColor Cyan
Write-Host "1. Did Edge browser open successfully?" -ForegroundColor White
Write-Host "2. Was CDP connection established?" -ForegroundColor White
Write-Host "3. Were fingerprint values spoofed correctly?" -ForegroundColor White
Write-Host "4. Did you see console messages in DevTools?" -ForegroundColor White
Write-Host "5. Any errors in the terminal?" -ForegroundColor White
Write-Host ""
Write-Host "💡 This approach should work even if antivirus blocks Rod!" -ForegroundColor Green