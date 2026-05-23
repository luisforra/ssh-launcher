# sshl — SSH Host Launcher

A Go program that parses your SSH config (with `Include` support), presents an interactive arrow-key navigable list with type-to-filter, and launches `ssh.exe` on host selection.

## Usage

```powershell
sshl              # interactive list → select → connect
sshl myserver     # connect directly to myserver
sshl --help       # list all discovered hosts
```

## Features

- Parses `%USERPROFILE%\.ssh\config` and all `Include`d files recursively
- Filtered host list — wildcards (`*`, `?`, negations) are hidden; only concrete aliases shown
- Skips `Match` blocks (conditional hosts)
- Arrow-key navigation with real-time case-insensitive filtering
- Falls back to a numbered list if raw terminal mode is unavailable
- `sshl <host>` for direct connection without the interactive menu
- Handles circular includes via a visited-set tracker

## Build

```powershell
go build -o sshl.exe .
```

Requires Go 1.24+ and `golang.org/x/term`.

## Requirements

- Windows 10/11 with [OpenSSH Client](https://learn.microsoft.com/en-us/windows-server/administration/openssh/openssh-install-windows) enabled
- `~/.ssh/config` (or files included via `Include`)

## License

Public domain. No rights reserved.

## Credits

Original developer: Luis Forra <luis.forra@gmail.com>
Built with [deepseek-v4-flash-free](https://opencode.ai) via opencode.
