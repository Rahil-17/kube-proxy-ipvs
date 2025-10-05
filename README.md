# Kube Proxy IPVS - IPVS Load Balancer

A simple, clean IPVS proxy that creates virtual services and backends from YAML configuration.

## Quick Start

1. **Install dependencies:**
   ```bash
   go get gopkg.in/yaml.v3
   ```

2. **Install ipvsadm (Linux only):**
   ```bash
   # Ubuntu/Debian
   sudo apt-get install ipvsadm

   # CentOS/RHEL
   sudo yum install ipvsadm
   ```

3. **Run with example config:**
   ```bash
   sudo go run ./cmd/miniproxy --config=config.yaml
   ```

4. **Verify IPVS configuration:**
   ```bash
   # Show current IPVS status using the program
   sudo go run ./cmd/miniproxy --status

   # Or use ipvsadm directly
   sudo ipvsadm -Ln
   ```

### Building for Linux from Other Platforms

If you're developing on macOS/Windows but want to build for Linux:

```bash
# Build for Linux from macOS/Windows
GOOS=linux GOARCH=amd64 go build -o miniproxy-linux ./cmd/miniproxy

# Then copy to Linux machine and run
sudo ./miniproxy-linux --config=config.yaml
```

## Configuration

Edit `config.yaml` to define your virtual service and backends:

```yaml
service:
  vip: "10.10.10.10"      # Virtual IP
  port: 80                 # Service port
  protocol: "TCP"          # TCP or UDP
  scheduler: "rr"          # Round-robin scheduler

backends:
  - ip: "192.168.0.2"      # Backend server IP
    port: 8080             # Backend port
  - ip: "192.168.0.3"
    port: 8080
```

## Project Structure

```
kube-proxy-ipvs/
├── cmd/miniproxy/main.go     # CLI entry point
├── pkg/config/config.go      # YAML configuration
├── pkg/ipvs/ipvs.go         # IPVS command wrapper
├── config.yaml              # Example configuration
└── go.mod                   # Go dependencies
```

## Requirements

- **Linux only** - This program uses Linux-specific IPVS functionality
- IPVS kernel module loaded (`sudo modprobe ip_vs`)
- `ipvsadm` command-line tool installed
- Root privileges (for IPVS operations)
- Go 1.21+

> **Note**: This program will not run on macOS, Windows, or other non-Linux systems. It includes build constraints to prevent compilation on unsupported platforms.

## Usage

### Apply Configuration
```bash
sudo go run ./cmd/miniproxy --config=config.yaml
```

### Show Current Status
```bash
sudo go run ./cmd/miniproxy --status
```

### Command Line Options
- `--config`: Path to YAML configuration file (default: `config.yaml`)
- `--status`: Show current IPVS status and exit (doesn't apply any changes)
- `--help`: Show help message

## Testing

See [TEST-GUIDE.md](TEST-GUIDE.md) for complete instructions on how to test round robin load balancing using curl commands.

### Example Output

When you apply a configuration, you'll see output like:
```
2025/10/05 22:20:00 Added backend 192.168.0.2:8080
2025/10/05 22:20:00 Added backend 192.168.0.3:8080
2025/10/05 22:20:00 Created virtual server 10.10.10.10:80 [rr]
2025/10/05 22:20:00 Successfully applied IPVS config from config.yaml
2025/10/05 22:20:00 Current IPVS configuration:

IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port Scheduler Flags
  -> RemoteAddress:Port           Forward Weight ActiveConn InActConn
TCP  10.10.10.10:80 rr
  -> 192.168.0.2:8080            Masq    1      0          0
  -> 192.168.0.3:8080            Masq    1      0          0
```

## Implementation

This project uses a simple command-based approach that calls `ipvsadm` directly, making it:
- **Lightweight**: No complex dependencies
- **Reliable**: Uses the standard Linux IPVS tools
- **Simple**: Easy to understand and debug
- **Verifiable**: Automatically shows status after applying changes
