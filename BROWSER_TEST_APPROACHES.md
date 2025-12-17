# Ghost Browser Test Approaches

## Problem
The original Rod library-based browser launcher is being blocked by Windows Defender/antivirus software, preventing fingerprint spoofing from working properly.

## Available Test Approaches

### 1. Rod-based Test (Original)
**File:** `test-browser-launcher.go`
**Script:** `run-browser-test.ps1`

**Pros:**
- Full-featured browser automation
- Complete fingerprint spoofing
- Mature library with good documentation

**Cons:**
- Blocked by antivirus software
- May require antivirus exceptions
- External dependency

**Status:** ⚠️ Blocked by antivirus

### 2. Simple Edge Launcher
**File:** `test-simple-browser.go`
**Script:** Built into `run-all-tests.ps1`

**Pros:**
- No external dependencies
- Always works (no antivirus issues)
- Verifies Edge can be launched

**Cons:**
- No fingerprint spoofing
- Limited functionality
- Manual testing only

**Status:** ✅ Works but limited

### 3. Alternative Approach (Startup Scripts)
**File:** `test-alternative-browser.go`
**Script:** `run-alternative-test.ps1`

**Pros:**
- Bypasses Rod library
- Uses Edge's built-in script injection
- No antivirus blocking

**Cons:**
- Limited injection capabilities
- Requires Edge-specific features
- May not work on all Edge versions

**Status:** 🔄 Experimental

### 4. Direct CDP Implementation
**File:** `test-cdp-direct.go`
**Script:** `run-cdp-test.ps1`

**Pros:**
- Pure Go implementation
- No external dependencies
- Full CDP control
- No antivirus issues
- Complete fingerprint spoofing

**Cons:**
- More complex implementation
- Requires CDP knowledge
- Custom HTTP client needed

**Status:** ✅ Recommended solution

## Test Instructions

### Quick Test (Recommended)
```powershell
.\run-cdp-test.ps1
```

### Complete Test Suite
```powershell
.\run-all-tests.ps1
```

### Individual Tests
```powershell
# Rod-based (may be blocked)
.\run-browser-test.ps1

# Alternative approach
.\run-alternative-test.ps1

# Direct CDP
.\run-cdp-test.ps1
```

## What to Check

For each test, verify:

1. **Edge Launch**: Did Edge browser open?
2. **Fingerprint Spoofing**: Are values changed on browserleaks.com?
3. **Console Messages**: Do you see `[Ghost Browser]` messages in DevTools?
4. **Terminal Errors**: Any errors in the PowerShell terminal?
5. **CDP Connection**: (For CDP tests) Was the connection successful?

### Expected Spoofed Values
Check these values on browserleaks.com/javascript:
- `navigator.hardwareConcurrency`: Should match generated profile
- `navigator.deviceMemory`: Should match generated profile  
- `navigator.platform`: Should match generated profile
- `navigator.webdriver`: Should be `false` or `undefined`

## Recommended Solution

**Use Direct CDP Implementation** (`test-cdp-direct.go`)

This approach:
- ✅ Bypasses antivirus blocking
- ✅ Provides full fingerprint spoofing
- ✅ Uses only standard Go libraries
- ✅ Communicates directly with Edge via CDP
- ✅ Can be integrated into main application

## Integration Plan

Once testing confirms the CDP direct approach works:

1. Replace Rod dependency in `internal/browser/browser.go`
2. Implement CDP client using the working test code
3. Update fingerprint injection to use CDP Runtime.evaluate
4. Test with main Ghost Browser application
5. Update documentation and build scripts

## Troubleshooting

### Edge Not Found
- Ensure Microsoft Edge is installed
- Check Edge installation paths in code
- Try running Edge manually first

### CDP Connection Failed
- Ensure Edge launched with `--remote-debugging-port=9222`
- Check if port 9222 is available
- Verify no firewall blocking localhost connections

### Spoofing Not Working
- Check DevTools Console for error messages
- Verify script injection succeeded
- Test on multiple websites (browserleaks.com, creepjs.com)

### Antivirus Blocking
- Add Ghost Browser folder to antivirus exceptions
- Use CDP direct approach instead of Rod
- Check Windows Defender real-time protection settings