# Backup Procedure

## Overview

This document describes the standard backup procedure for production servers.

## Pre-requisites

- SSH access to the target server
- Sufficient disk space on backup storage
- Backup user credentials

## Steps

### 1. Verify Current State

```bash
# Check disk space
df -h

# Check running services
systemctl list-units --type=service --state=running
```

### 2. Create Database Backup

```bash
# PostgreSQL
pg_dump -U postgres dbname > backup_$(date +%Y%m%d).sql

# MySQL
mysqldump -u root -p dbname > backup_$(date +%Y%m%d).sql
```

### 3. Backup Configuration Files

```bash
tar -czvf config_backup_$(date +%Y%m%d).tar.gz /etc/nginx /etc/postgresql
```

### 4. Verify Backup Integrity

```bash
# Check tar archive
tar -tzvf config_backup_*.tar.gz

# Check SQL dump
head -20 backup_*.sql
```

## Rollback Procedure

If backup fails:

1. Check disk space
2. Verify database connectivity
3. Check permissions
4. Review logs: `journalctl -xe`

## Related Commands

- `filesystem/disk-usage` - Check available space
- `services/status` - Verify service status
