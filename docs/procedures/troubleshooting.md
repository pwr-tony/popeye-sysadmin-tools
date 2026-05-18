# Troubleshooting Guide

## Common Issues

### High CPU Usage

1. Identify the process:
```bash
top -o %CPU
ps aux --sort=-%cpu | head -10
```

2. Check if it's a known service
3. Review recent changes

### High Memory Usage

1. Check memory status:
```bash
free -h
```

2. Find memory-hungry processes:
```bash
ps aux --sort=-%mem | head -10
```

### Disk Space Issues

1. Check usage:
```bash
df -h
```

2. Find large files:
```bash
du -ah / | sort -rh | head -20
```

3. Clean up:
- Old logs: `/var/log`
- Package cache: `dnf clean all` or `apt clean`
- Old kernels

### Network Connectivity

1. Check interface status:
```bash
ip addr show
```

2. Test connectivity:
```bash
ping 8.8.8.8
```

3. Check DNS:
```bash
dig google.com
```

4. Check routes:
```bash
ip route show
```

### Service Not Starting

1. Check status:
```bash
systemctl status service_name
```

2. Check logs:
```bash
journalctl -u service_name -n 100
```

3. Check configuration:
```bash
systemctl cat service_name
```

## Log Locations

| Service | Log Location |
|---------|-------------|
| System | `/var/log/messages` or `journalctl` |
| SSH | `/var/log/secure` |
| Web (nginx) | `/var/log/nginx/` |
| Database | `/var/log/postgresql/` |

## Related Commands

- `services/logs` - View service logs
- `network/show-interfaces` - Check network
- `filesystem/disk-usage` - Check disk space
