# sshl — SSH Host Launcher

## Overview
A Go program that parses SSH config (with `Include` support), presents an interactive
arrow-key navigable list with type-to-filter, and launches `ssh.exe` on host selection.

## File structure
```
sshl/
├── go.mod          — module: sshl
├── main.go         — entry point, orchestration, CLI arg handling
├── config.go       — SSH config parser (recursive Include, cycle detection)
├── tui.go          — interactive terminal UI (golang.org/x/term)
└── AGENTS.md       — this file
```

## Build & run
```powershell
# Build (output .exe)
go build -o sshl.exe .

# Or install to $GOPATH/bin
go install .

# Run
.\sshl.exe
.\sshl.exe myserver           # direct connect
```

## Features
- Parses `%USERPROFILE%\.ssh\config` and all `Include`d files
- Filters wildcard hosts (`*`, `?`, negations) — only shows concrete aliases
- Skips `Match` blocks
- Handles circular includes (visited set)
- Arrow-key navigation with type-to-filter (case-insensitive substring)
- Fallback to numbered list if raw terminal mode unavailable
- `sshl <host>` for direct connection without interactive menu
- `sshl --help` lists all discovered hosts

## Dependencies
- `golang.org/x/term` — raw terminal I/O (Windows + Unix)

## Testing
```powershell
# build only
go build ./...

# vet
go vet ./...
```

## Key types
```go
type HostEntry struct {
    Alias        string  // the SSH alias (e.g. "myserver")
    HostName     string  // resolved HostName if set, else alias
    User         string  // SSH user
    Port         string  // SSH port
    IdentityFile string  // path to identity file
    FileName     string  // source config file path
}

## License
Public domain. No rights reserved.

## Credits
- Original developer: Luis Forra <luis.forra@gmail.com>
- Built with [deepseek-v4-flash-free](https://opencode.ai) via opencode
```
