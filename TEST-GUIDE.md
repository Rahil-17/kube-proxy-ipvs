# IPVS Round Robin Testing Guide

This guide will help you validate that round robin load balancing is working correctly using curl.

## Prerequisites

- Linux system with IPVS support
- `ipvsadm` installed
- `curl` installed
- Root/sudo access

## Test Setup

### Step 1: Start Backend Servers

Run in **Terminal 1**:
```bash
./setup-test-backends.sh
```

You should see:
```
Starting server on port 9001...
Starting server on port 9002...
Starting server on port 9003...
All test servers started!
```

### Step 2: Verify Backends Work

In **Terminal 2**, test each backend directly:
```bash
curl http://127.0.0.1:9001
# Should show: Backend Server 1 (Port 9001)

curl http://127.0.0.1:9002
# Should show: Backend Server 2 (Port 9002)

curl http://127.0.0.1:9003
# Should show: Backend Server 3 (Port 9003)
```

## Apply IPVS Configuration

In **Terminal 2**:
```bash
sudo go run ./cmd/miniproxy --config=config.yaml
```

Expected output:
```
Added backend 127.0.0.1:9001
Added backend 127.0.0.1:9002
Added backend 127.0.0.1:9003
Created virtual server 127.0.0.1:8080 [rr]
Successfully applied IPVS config from config.yaml
Current IPVS configuration:
...
```

## Test Round Robin Routing

### Method 1: Sequential Curl Commands

Run these commands one by one and observe the responses:

```bash
curl http://127.0.0.1:8080
# Expected: Backend Server 1 (Port 9001)

curl http://127.0.0.1:8080
# Expected: Backend Server 2 (Port 9002)

curl http://127.0.0.1:8080
# Expected: Backend Server 3 (Port 9003)

curl http://127.0.0.1:8080
# Expected: Backend Server 1 (Port 9001) - Cycle repeats

curl http://127.0.0.1:8080
# Expected: Backend Server 2 (Port 9002)

curl http://127.0.0.1:8080
# Expected: Backend Server 3 (Port 9003)
```

### Method 2: Loop Test

Run this command to test multiple times automatically:

```bash
for i in {1..12}; do
  echo "Request $i:"
  curl -s http://127.0.0.1:8080 | grep -oP '(?<=Backend Server )\d+' || echo "Failed"
  echo ""
done
```

Expected output shows rotation through servers 1, 2, 3, 1, 2, 3, etc.

### Method 3: Check IPVS Statistics

Monitor the connection distribution:

```bash
# Check IPVS stats before testing
sudo ipvsadm -Ln --stats

# Run some curl commands
for i in {1..9}; do curl -s http://127.0.0.1:8080 > /dev/null; done

# Check stats again - should show equal distribution
sudo ipvsadm -Ln --stats
```

You should see roughly equal connection counts across all three backends:
```
TCP  127.0.0.1:8080 rr
  -> 127.0.0.1:9001            Masq    3      0          0
  -> 127.0.0.1:9002            Masq    3      0          0
  -> 127.0.0.1:9003            Masq    3      0          0
```

## Cleanup

### Stop Backend Servers
In Terminal 1 (where servers are running):
```bash
# Press Ctrl+C
```

### Clear IPVS Rules
```bash
sudo ipvsadm -C
```

## Troubleshooting

### Issue: "Connection refused"
- Ensure backend servers are running
- Check with: `netstat -tulpn | grep -E '9001|9002|9003'`

### Issue: "No response from VIP"
- Verify IPVS configuration: `sudo ipvsadm -Ln`
- Check if IPVS module is loaded: `lsmod | grep ip_vs`

### Issue: "Always hitting same backend"
- May need to close connections between requests
- Try: `curl --no-keepalive http://127.0.0.1:8080`

### Issue: macOS/Windows
- This test requires Linux. Use a Linux VM or container.
