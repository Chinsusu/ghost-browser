# Ghost Browser - Advanced Spoofing v4
# Features:
# - Random fingerprint generation
# - CSS Media Query protection
# - WebRTC complete block
# - Canvas/Audio noise with consistent seed
# - Timezone sync

Write-Host "=== Ghost Browser v4 - Advanced Spoofing ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "Features:" -ForegroundColor Yellow
Write-Host "  [+] Random fingerprint per launch"
Write-Host "  [+] CSS Media Query protection"
Write-Host "  [+] WebRTC complete disable"
Write-Host "  [+] Canvas/Audio noise (seeded)"
Write-Host "  [+] Battery API spoof"
Write-Host "  [+] Permissions API spoof"
Write-Host ""

go mod tidy 2>$null
go run main.go
