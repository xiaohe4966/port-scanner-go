# Port Scanner 🔍

A simple and fast LAN port scanner written in Go.

## Files

| File | Description |
|------|-------------|
| `main.go` | **Single-port scanner** — Scans a specific port (default `8000`) across all hosts in a `/24` subnet (e.g. `192.168.1.1–254`). Useful for quickly finding which machines in your LAN have a particular service running. |
| `all.go` | **All-ports scanner** — Scans **all 65535 TCP ports** on every host in a `/24` subnet. Shows real-time progress and outputs a per-host summary of open ports. Useful for comprehensive network auditing. |

## Features

- 🚀 Concurrent scanning with configurable worker count (default 100 goroutines)
- ⏱️ Adjustable connection timeout
- 📊 Real-time progress display (`all.go`)
- 📋 Clean summary of open ports

## Usage

### Single-port scan (`main.go`)

```bash
go run main.go
```

Scans `192.168.1.1–254` on port `8000`. Edit the constants at the top of `main.go` to change:

- `subnet` — target subnet prefix
- `startIP` / `endIP` — IP range
- `port` — target port
- `timeout` — connection timeout
- `workers` — concurrency level

### All-ports scan (`all.go`)

```bash
go run all.go
```

Scans `192.168.1.1–254` on **all 65535 ports**. Edit the constants at the top of `all.go` to change:

- `subnet` — target subnet prefix
- `startIP` / `endIP` — IP range
- `timeout` — connection timeout
- `workers` — concurrency level

> ⚠️ Scanning all 65535 ports on 254 hosts = **~16.5M connections**. Make sure your network and system can handle it.

## License

MIT
