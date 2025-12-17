# Ghost Browser Integration Guide

## Files to Replace/Add

Replace these files in your ghost-browser project:

### 1. Replace `internal/browser/browser.go`
Copy the new browser.go that uses chromedp instead of rod.

### 2. Replace `internal/browser/spoof.go`  
Copy the new spoof.go with advanced fingerprint spoofing.

### 3. Replace `internal/fingerprint/fingerprint.go`
Copy the updated fingerprint struct.

### 4. Replace `internal/fingerprint/generator.go`
Copy the improved generator with more realistic values.

### 5. Update `go.mod`
Add chromedp dependencies:
```
github.com/chromedp/cdproto v0.0.0-20241022234722-4571571f02d1
github.com/chromedp/chromedp v0.11.0
```

## Installation Steps

1. **Backup your current files**
```powershell
cd ghost-browser
mkdir backup
copy internal\browser\*.go backup\
copy internal\fingerprint\*.go backup\
```

2. **Copy new files**
```powershell
# Copy from the integration folder
copy path\to\ghost-browser-integration\internal\browser\*.go internal\browser\
copy path\to\ghost-browser-integration\internal\fingerprint\*.go internal\fingerprint\
```

3. **Update go.mod**
Add the chromedp imports to your go.mod:
```
go get github.com/chromedp/chromedp@v0.11.0
go get github.com/chromedp/cdproto@latest
```

4. **Remove rod dependency** (optional)
```
go mod edit -droprequire github.com/go-rod/rod
```

5. **Tidy up**
```powershell
go mod tidy
```

6. **Test**
```powershell
wails dev
```

## What Changed

### browser.go
- Changed from `rod` to `chromedp`
- Uses `page.AddScriptToEvaluateOnNewDocument` for proper injection
- WebRTC disabled via Chrome flags
- Better context management

### spoof.go
- Complete fingerprint spoofing script
- WebRTC complete block
- Canvas/Audio noise with seeded randomness
- CSS Media Query protection
- Battery/Permissions API spoof
- Automation flag removal

### fingerprint/generator.go
- More realistic GPU list
- Timezone + Locale sync
- Better language matching
- User agent generation

## Testing

After integration, test with:
1. Launch the app: `wails dev`
2. Create a profile
3. Click "Launch Browser"
4. Browser should open with spoofed fingerprint
5. Visit https://browserleaks.com/javascript to verify
6. Visit https://abrahamjuliot.github.io/creepjs/ for full test

## Expected Results

- hardwareConcurrency: Random (4/6/8/12/16)
- deviceMemory: Random (4/8/16/32)
- Screen: Random common resolution
- WebRTC: Blocked (no IP leak)
- Canvas: Noise added
- Headless: Blocked
