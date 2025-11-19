#!/bin/bash

# Automated backup script for cron jobs
# Add to crontab: 0 2 * * * /path/to/wiki-notes/scripts/auto-backup.sh

cd "$(dirname "$0")/.."
./scripts/backup-db.sh

# Optional: Sync to cloud storage (uncomment and configure)
# Example with rclone:
# rclone copy ./backups remote:wiki-notes-backups

# Example with AWS S3:
# aws s3 sync ./backups s3://your-bucket/wiki-notes-backups/

# Example with rsync to remote server:
# rsync -avz ./backups/ user@backup-server:/backups/wiki-notes/

