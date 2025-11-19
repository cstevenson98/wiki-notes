# Database Backup & Recovery Guide

## Overview

This guide explains how to backup and restore your Wiki Notes database to prevent accidental data loss.

## Quick Start

### Manual Backup

```bash
# Create a backup
./scripts/backup-db.sh
```

Backups are saved to `./backups/` with timestamps like: `wiki_notes_backup_20231119_143022.sql.gz`

### Restore from Backup

```bash
# List available backups
ls -lh backups/

# Restore a specific backup
./scripts/restore-db.sh backups/wiki_notes_backup_20231119_143022.sql.gz
```

## Protection Strategies

### 1. Prevent Accidental Volume Deletion

The docker-compose volume is configured with a named volume that persists even if you run `docker-compose down`.

⚠️ **WARNING**: `docker-compose down -v` will still delete volumes. Avoid using the `-v` flag unless you intend to wipe data.

To make the volume completely external (maximum protection):

```bash
# Create the volume manually
docker volume create wiki_notes_db_data

# Edit docker-compose.yml and set external: true
# volumes:
#   postgres_data:
#     name: wiki_notes_db_data
#     external: true
```

### 2. Automated Backups

#### Set up daily backups with cron:

```bash
# Edit crontab
crontab -e

# Add this line for daily backups at 2 AM:
0 2 * * * cd /home/conor/dev/wiki-notes && ./scripts/auto-backup.sh >> logs/backup.log 2>&1

# Or for hourly backups:
0 * * * * cd /home/conor/dev/wiki-notes && ./scripts/auto-backup.sh >> logs/backup.log 2>&1
```

#### Backup retention:

Edit `scripts/backup-db.sh` and uncomment these lines to keep only the last N backups:

```bash
MAX_BACKUPS=10
cd "${BACKUP_DIR}"
ls -t wiki_notes_backup_*.sql.gz | tail -n +$((MAX_BACKUPS + 1)) | xargs -r rm
```

### 3. Off-Site Backups

#### Option A: Cloud Storage (Recommended)

**Using rclone (supports S3, Google Drive, Dropbox, etc.):**

```bash
# Install rclone
curl https://rclone.org/install.sh | sudo bash

# Configure remote storage
rclone config

# Edit scripts/auto-backup.sh and uncomment:
rclone copy ./backups remote:wiki-notes-backups
```

**Using AWS S3:**

```bash
# Install AWS CLI
sudo apt install awscli

# Configure credentials
aws configure

# Edit scripts/auto-backup.sh and uncomment:
aws s3 sync ./backups s3://your-bucket/wiki-notes-backups/
```

#### Option B: Remote Server Backup

```bash
# Set up SSH key authentication to backup server
ssh-copy-id user@backup-server

# Edit scripts/auto-backup.sh and uncomment:
rsync -avz ./backups/ user@backup-server:/backups/wiki-notes/
```

#### Option C: Git-based Backup

```bash
# Initialize git in backups directory
cd backups
git init
git add .
git commit -m "Backup $(date)"

# Push to private remote repository
git remote add origin git@github.com:yourusername/wiki-notes-backups.git
git push origin main
```

### 4. Docker Volume Backup

Alternatively, backup the entire Docker volume:

```bash
# Backup volume to tar file
docker run --rm \
  -v wiki_notes_db_data:/data \
  -v $(pwd)/backups:/backup \
  alpine tar czf /backup/volume_backup_$(date +%Y%m%d_%H%M%S).tar.gz -C /data .

# Restore volume from tar file
docker run --rm \
  -v wiki_notes_db_data:/data \
  -v $(pwd)/backups:/backup \
  alpine tar xzf /backup/volume_backup_20231119_143022.tar.gz -C /data
```

## Backup Best Practices

1. **3-2-1 Rule**: Keep 3 copies of data, on 2 different media, with 1 copy off-site
2. **Test Restores**: Regularly test that backups can be restored
3. **Automate**: Set up automated backups with cron
4. **Monitor**: Check backup logs regularly
5. **Retention**: Keep multiple backup versions (daily for 7 days, weekly for 4 weeks, monthly for 12 months)

## Disaster Recovery Procedure

If you lose your database:

1. **Stop the application:**
   ```bash
   docker-compose down
   ```

2. **Find the most recent backup:**
   ```bash
   ls -lht backups/ | head -5
   ```

3. **Restore the database:**
   ```bash
   docker-compose up -d db
   # Wait for database to be ready
   ./scripts/restore-db.sh backups/wiki_notes_backup_YYYYMMDD_HHMMSS.sql.gz
   ```

4. **Restart the application:**
   ```bash
   docker-compose up -d
   ```

## Monitoring Backup Health

Create a simple script to check backup status:

```bash
#!/bin/bash
# scripts/check-backups.sh

BACKUP_DIR="./backups"
LATEST_BACKUP=$(ls -t ${BACKUP_DIR}/wiki_notes_backup_*.sql.gz 2>/dev/null | head -1)

if [ -z "$LATEST_BACKUP" ]; then
    echo "❌ No backups found!"
    exit 1
fi

BACKUP_AGE=$((($(date +%s) - $(stat -c %Y "$LATEST_BACKUP")) / 86400))

if [ $BACKUP_AGE -gt 1 ]; then
    echo "⚠️  Latest backup is $BACKUP_AGE days old: $LATEST_BACKUP"
    exit 1
else
    echo "✅ Latest backup: $LATEST_BACKUP (${BACKUP_AGE} days old)"
fi
```

## Backup File Locations

- **Local backups**: `./backups/`
- **Backup logs**: `./logs/backup.log` (if using cron)
- **Docker volume**: `wiki_notes_db_data`

## Security Notes

- Backups contain all your wiki data in plain text
- Encrypt backups if storing off-site:
  ```bash
  gpg -c backups/wiki_notes_backup_YYYYMMDD.sql.gz
  ```
- Don't commit backups to public git repositories
- Set appropriate file permissions: `chmod 600 backups/*.sql.gz`

## Troubleshooting

### Backup script fails with "container not running"

```bash
# Check if database is running
docker-compose ps

# Start only the database
docker-compose up -d db
```

### Backup file is empty or corrupt

```bash
# Test the backup
gunzip -c backups/wiki_notes_backup_YYYYMMDD.sql.gz | head -20
```

### Out of disk space

```bash
# Check disk usage
df -h

# Clean old backups manually
rm backups/wiki_notes_backup_2023*.sql.gz
```

## Additional Resources

- [PostgreSQL Backup Documentation](https://www.postgresql.org/docs/current/backup.html)
- [Docker Volume Management](https://docs.docker.com/storage/volumes/)
- [rclone Documentation](https://rclone.org/docs/)

