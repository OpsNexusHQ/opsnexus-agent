# OpsNexus Agent (`opsnexus-agent`)

[![Release](https://img.shields.io/badge/release-v0.5.0-blue.svg)](https://github.com/OpsNexusHQ/opsnexus-agent/releases/tag/v0.5.0)
[![Go Version](https://img.shields.io/badge/go-1.25+-00ADD8.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Lightweight, high-performance Linux system monitoring agent for **OpsNexus**. Written in Go and powered by `github.com/shirou/gopsutil/v4`, it collects real-time hardware and OS metrics with minimal resource usage (< 15MB RAM).

---

## ✨ Implemented Capabilities (v0.5.0)

- **CPU Metrics**: Usage percentage, per-core stats, physical core counts, model info.
- **Memory Metrics**: Total, used, free, available RAM, swap usage percentage.
- **Disk Metrics**: Partition mount points, total/used/free space, disk usage percentage.
- **Network Metrics**: Bytes sent/received, packets, interface status, errors.
- **System Uptime & Host Details**: System boot time, uptime in seconds, OS/arch metadata.
- **Running Processes**: Active system process counts.
- **Automatic Registration**: Self-registers with the OpsNexus backend on startup.
- **Graceful Lifecycle**: Supports context cancellation and clean shutdown signals (`SIGINT`, `SIGTERM`).

---

## ⚙️ Configuration & Environment Variables

The agent can be configured via environment variables or flag parameters:

| Variable | Default | Description |
|---|---|---|
| `OPSNEXUS_BACKEND_URL` | `http://localhost:8080` | URL of the OpsNexus backend server |
| `OPSNEXUS_COLLECT_INTERVAL` | `10s` | Metric collection & transmission frequency |
| `OPSNEXUS_AGENT_NAME` | Hostname | Human-readable name for the agent |

---

## 🚀 Quickstart & Installation

### Option 1: Run from Source
```bash
# 1. Clone
git clone https://github.com/OpsNexusHQ/opsnexus-agent.git
cd opsnexus-agent

# 2. Configure & Run
export OPSNEXUS_BACKEND_URL="http://localhost:8080"
export OPSNEXUS_COLLECT_INTERVAL="10s"

go run ./cmd/opsnexus-agent
```

### Option 2: Build Standalone Binary
```bash
go build -o opsnexus-agent ./cmd/opsnexus-agent
./opsnexus-agent
```

---

## 🏛️ System Architecture Placement

```text
┌────────────────────────────────────────────────────────┐
│                      Linux Host                        │
│  ┌──────────────────────────────────────────────────┐  │
│  │                 opsnexus-agent                   │  │
│  │   gopsutil/v4 CPU/RAM/Disk/Network Collector    │  │
│  └─────────────────────────┬────────────────────────┘  │
└────────────────────────────┼───────────────────────────┘
                             │
                  HTTP POST /telemetry (10s)
                             │
                             ▼
                    ┌─────────────────┐
                    │ OpsNexus Backend│
                    └─────────────────┘
```

---

## 🗺️ Roadmap (Future Scope)

- [ ] **Docker & Container Metrics**: Container CPU/RAM utilization monitoring.
- [ ] **eBPF Tracing**: Low-overhead kernel network flow analysis.
- [ ] **Log Forwarding**: Structured log tailing and backend ingestion.
- [ ] **Windows & macOS Support**: Cross-platform system collection.

---

## 📄 License

Part of the [OpsNexus](https://github.com/OpsNexusHQ) ecosystem. Licensed under the MIT License.
