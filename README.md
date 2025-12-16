# Ghost Browser

An antidetect browser built with Go, Wails, and React. Uses Microsoft Edge via Chrome DevTools Protocol (CDP) for browser automation with advanced fingerprint spoofing.

## Features

- **Profile Management**: Create, edit, duplicate, import/export browser profiles
- **Fingerprint Spoofing**: Canvas, WebGL, Audio, Navigator, Screen, Timezone, WebRTC
- **Proxy Support**: HTTP, HTTPS, SOCKS4, SOCKS5 with health checking
- **AI Personalities**: Local LLM integration via Ollama for unique personality per profile
- **Scheduling**: Define activity schedules and browsing patterns per profile
- **Offline**: All data stored locally, no cloud dependencies

## Requirements

- Windows 10/11
- Microsoft Edge (pre-installed on Windows)
- Go 1.22+
- Node.js 18+
- Wails CLI v2
- (Optional) Ollama for AI features

## Installation

### 1. Install Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 2. Install Dependencies

```bash
# Go dependencies
go mod tidy

# Frontend dependencies
cd frontend
npm install
cd ..
```

### 3. Development Mode

```bash
wails dev
```

### 4. Build for Production

```bash
wails build
```

The executable will be in `build/bin/`.

## Project Structure

```
ghost-browser/
├── cmd/ghost/          # Application entry point
├── internal/
│   ├── app/           # Main application controller
│   ├── browser/       # Browser launcher and spoofing
│   ├── database/      # SQLite database layer
│   ├── fingerprint/   # Fingerprint generation
│   ├── profile/       # Profile management
│   ├── proxy/         # Proxy management
│   └── ai/            # AI personality engine
├── frontend/          # React frontend
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   └── styles/
│   └── ...
├── configs/           # Configuration files
├── scripts/           # Build scripts
└── docs/             # Documentation
```

## Fingerprint Spoofing

The browser spoofs the following fingerprints:

| Fingerprint | Method |
|------------|--------|
| User Agent | Navigator override |
| Screen Resolution | Screen prototype override |
| Hardware Concurrency | Navigator override |
| Device Memory | Navigator override |
| WebGL Vendor/Renderer | getParameter intercept |
| Canvas | Noise injection |
| Audio | AudioContext noise |
| Timezone | Date/Intl override |
| WebRTC | Complete disable or IP masking |
| Plugins | PluginArray spoof |

## AI Personalities (Ollama)

To use AI personalities:

1. Install Ollama: https://ollama.ai
2. Pull a model: `ollama pull llama3.2`
3. Start Ollama: `ollama serve`

Each profile can have a unique AI personality with:
- Name, age, gender, occupation
- Interests and expertise
- Writing style (formal/casual, emojis, etc.)
- Typing speed and patterns
- Activity schedule

## Configuration

Copy `configs/config.example.yaml` to `configs/config.yaml` and customize:

```yaml
browser:
  executable_path: ""  # Auto-detect Edge

ai:
  ollama_url: "http://localhost:11434"
  model: "llama3.2"
```

## API Reference

The app exposes these methods to the frontend via Wails bindings:

### Profiles
- `GetProfiles()` - List all profiles
- `GetProfile(id)` - Get profile by ID
- `CreateProfile(name, options)` - Create new profile
- `UpdateProfile(profile)` - Update profile
- `DeleteProfile(id)` - Delete profile
- `GenerateRandomProfile()` - Generate random profile
- `DuplicateProfile(id)` - Duplicate profile
- `ExportProfile(id, path)` - Export to file
- `ImportProfile(path)` - Import from file

### Proxies
- `GetProxies()` - List all proxies
- `AddProxy(proxy)` - Add proxy
- `DeleteProxy(id)` - Delete proxy
- `CheckProxy(id)` - Check proxy health
- `CheckAllProxies()` - Check all proxies
- `ImportProxies(text, format)` - Bulk import

### Browser
- `LaunchBrowser(profileId)` - Launch browser
- `CloseBrowser(profileId)` - Close browser
- `GetRunningBrowsers()` - List running browsers

### AI
- `GetPersonality(profileId)` - Get AI personality
- `UpdatePersonality(profileId, personality)` - Update personality
- `GeneratePersonality()` - Generate random personality
- `Chat(profileId, message)` - Chat with personality
- `GetSchedule(profileId)` - Get schedule
- `UpdateSchedule(profileId, schedule)` - Update schedule

## License

MIT

## Disclaimer

This tool is for legitimate use cases such as:
- Social media management
- E-commerce operations
- Web development testing
- Privacy protection

Do not use for fraud, spam, or any illegal activities.
