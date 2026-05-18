# Popeye - Sysadmin Tools Manager

TUI-based tool manager for system administrators. Execute predefined commands, manage a local CMDB, and render procedures in Markdown.

## Features

- **Commands**: Execute sysadmin commands with parameters (firewall, network, filesystem, services, users)
- **CMDB**: Local infrastructure inventory (servers, networks)
- **Docs**: Markdown documentation viewer for procedures

## Installation

```bash
# Clone repository
git clone https://github.com/user/popeye.git
cd popeye

# Install dependencies
make deps

# Build
make build

# Run
./bin/popeye
```

## Usage

### TUI Mode

```bash
popeye
```

### Keybindings

| Key | Action |
|-----|--------|
| `q` / `Ctrl+C` | Quit |
| `?` | Help |
| `1` | Commands view |
| `2` | CMDB view |
| `3` | Docs view |
| `Enter` | Execute/Select |
| `Esc` | Back |
| `/` | Search |
| `↑/k` `↓/j` | Navigate |

## Command Categories

- **Firewall**: add-port, remove-port, list-ports, reload, add-service
- **Network**: show-interfaces, show-routes, show-connections, dns-check, ping
- **Filesystem**: disk-usage, mount-list, inode-usage, largest-files
- **Services**: status, start, stop, restart, enable, disable, logs
- **Users**: list-users, last-logins, who, failed-logins, user-info

## Configuration

### Adding Commands

Create YAML files in `configs/commands/`:

```yaml
category: mycategory
icon: "🔧"
commands:
  - name: my-command
    description: Description here
    params:
      - name: param1
        type: string
        required: true
        prompt: "Enter value"
    command: echo {{.param1}}
    sudo: false
```

### CMDB

Edit `data/cmdb.json` to add servers and networks:

```json
{
  "servers": [
    {
      "id": "srv-001",
      "hostname": "server-01",
      "ip": "192.168.1.10",
      "os": "Rocky Linux 9",
      "environment": "production",
      "tags": ["web"],
      "notes": "Web server"
    }
  ],
  "networks": []
}
```

### Documentation

Add Markdown files to `docs/procedures/`.

## Project Structure

```
popeye/
├── cmd/popeye/          # Entry point
├── internal/
│   ├── tui/             # TUI components
│   ├── commands/        # Command loader and executor
│   ├── cmdb/            # CMDB storage
│   └── docs/            # Markdown loader/renderer
├── configs/commands/    # YAML command definitions
├── data/                # CMDB data
└── docs/procedures/     # Markdown documentation
```

## License

MIT
