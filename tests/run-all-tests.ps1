#!/usr/bin/env pwsh

Write-Host "🧪 === Ghost Browser Complete Test Suite ===" -ForegroundColor Green
Write-Host ""

Write-Host "This will run all available browser launcher tests:" -ForegroundColor Yellow
Write-Host "  1. Rod-based test (may be blocked by antivirus)" -ForegroundColor Gray
Write-Host "  2. Simple Edge launcher (no spoofing)" -ForegroundColor Gray
Write-Host "  3. Alternative approach (startup scripts)" -ForegroundColor Gray
Write-Host "  4. Direct CDP implementation (pure Go)" -ForegroundColor Gray
Write-Host ""

$testResults = @()

# Test 1: Rod-based (original)
Write-Host "🔬 TEST 1: Rod-based Browser Launcher" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
try {
    go build -o test-browser-launcher.exe test-browser-launcher.go
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Rod test built successfully" -ForegroundColor Green
        Write-Host "⚠️  Note: This may be blocked by antivirus" -ForegroundColor Yellow
        Write-Host "Running Rod test..." -ForegroundColor White
        .\test-browser-launcher.exe
        $testResults += "Rod Test: Executed (check for antivirus blocking)"
        Remove-Item -Path "test-browser-launcher.exe" -ErrorAction SilentlyContinue
    } else {
        Write-Host "❌ Rod test build failed" -ForegroundColor Red
        $testResults += "Rod Test: Build Failed"
    }
} catch {
    Write-Host "❌ Rod test error: $_" -ForegroundColor Red
    $testResults += "Rod Test: Error - $_"
}

Write-Host "`n" -NoNewline
Read-Host "Press Enter to continue to next test"

# Test 2: Simple Edge launcher
Write-Host "`n🔬 TEST 2: Simple Edge Launcher (No Spoofing)" -ForegroundColor Cyan
Write-Host "===============================================" -ForegroundColor Cyan
try {
    go build -o test-simple-browser.exe test-simple-browser.go
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Simple test built successfully" -ForegroundColor Green
        Write-Host "Running simple Edge test..." -ForegroundColor White
        .\test-simple-browser.exe
        $testResults += "Simple Test: Executed"
        Remove-Item -Path "test-simple-browser.exe" -ErrorAction SilentlyContinue
    } else {
        Write-Host "❌ Simple test build failed" -ForegroundColor Red
        $testResults += "Simple Test: Build Failed"
    }
} catch {
    Write-Host "❌ Simple test error: $_" -ForegroundColor Red
    $testResults += "Simple Test: Error - $_"
}

Write-Host "`n" -NoNewline
Read-Host "Press Enter to continue to next test"

# Test 3: Alternative approach
Write-Host "`n🔬 TEST 3: Alternative Approach (Startup Scripts)" -ForegroundColor Cyan
Write-Host "=================================================" -ForegroundColor Cyan
try {
    go build -o test-alternative-browser.exe test-alternative-browser.go
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Alternative test built successfully" -ForegroundColor Green
        Write-Host "Running alternative approach test..." -ForegroundColor White
        .\test-alternative-browser.exe
        $testResults += "Alternative Test: Executed"
        Remove-Item -Path "test-alternative-browser.exe" -ErrorAction SilentlyContinue
    } else {
        Write-Host "❌ Alternative test build failed" -ForegroundColor Red
        $testResults += "Alternative Test: Build Failed"
    }
} catch {
    Write-Host "❌ Alternative test error: $_" -ForegroundColor Red
    $testResults += "Alternative Test: Error - $_"
}

Write-Host "`n" -NoNewline
Read-Host "Press Enter to continue to final test"

# Test 4: Direct CDP
Write-Host "`n🔬 TEST 4: Direct CDP Implementation (Pure Go)" -ForegroundColor Cyan
Write-Host "===============================================" -ForegroundColor Cyan
try {
    go build -o test-cdp-direct.exe test-cdp-direct.go
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ CDP direct test built successfully" -ForegroundColor Green
        Write-Host "Running direct CDP test..." -ForegroundColor White
        .\test-cdp-direct.exe
        $testResults += "CDP Direct Test: Executed"
        Remove-Item -Path "test-cdp-direct.exe" -ErrorAction SilentlyContinue
    } else {
        Write-Host "❌ CDP direct test build failed" -ForegroundColor Red
        $testResults += "CDP Direct Test: Build Failed"
    }
} catch {
    Write-Host "❌ CDP direct test error: $_" -ForegroundColor Red
    $testResults += "CDP Direct Test: Error - $_"
}

# Summary
Write-Host "`n🏁 === TEST SUITE SUMMARY ===" -ForegroundColor Green
Write-Host "=============================" -ForegroundColor Green
foreach ($result in $testResults) {
    Write-Host "  • $result" -ForegroundColor White
}

Write-Host "`n📊 RECOMMENDED NEXT STEPS:" -ForegroundColor Cyan
Write-Host "1. If Rod test worked: Use original implementation" -ForegroundColor White
Write-Host "2. If Rod was blocked: Use CDP Direct approach" -ForegroundColor White
Write-Host "3. If CDP failed: Use Simple launcher + manual injection" -ForegroundColor White
Write-Host "4. Report which test worked best for integration" -ForegroundColor White

Write-Host "`n✅ Complete test suite finished!" -ForegroundColor Green