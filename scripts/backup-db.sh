#!/bin/bash

# Database backup script for Wiki Notes
# Creates timestamped SQL dumps of the PostgreSQL database

set -e

# Configuration
BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="wiki_notes_backup_${TIMESTAMP}.sql"
CONTAINER_NAME="wiki-notes-db-1"
DB_NAME="wikinotes"
DB_USER="postgres"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Create backup directory if it doesn't exist
mkdir -p "${BACKUP_DIR}"

echo -e "${YELLOW}Starting database backup...${NC}"

# Check if container is running
if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo -e "${RED}Error: Database container '${CONTAINER_NAME}' is not running${NC}"
    exit 1
fi

# Create backup
echo -e "${YELLOW}Creating backup: ${BACKUP_FILE}${NC}"
docker exec -t "${CONTAINER_NAME}" pg_dump -U "${DB_USER}" "${DB_NAME}" > "${BACKUP_DIR}/${BACKUP_FILE}"

# Compress backup
echo -e "${YELLOW}Compressing backup...${NC}"
gzip "${BACKUP_DIR}/${BACKUP_FILE}"
BACKUP_FILE="${BACKUP_FILE}.gz"

# Get file size
BACKUP_SIZE=$(du -h "${BACKUP_DIR}/${BACKUP_FILE}" | cut -f1)

echo -e "${GREEN}✓ Backup completed successfully!${NC}"
echo -e "${GREEN}  File: ${BACKUP_DIR}/${BACKUP_FILE}${NC}"
echo -e "${GREEN}  Size: ${BACKUP_SIZE}${NC}"

# Optional: Keep only last N backups (uncomment to enable)
# MAX_BACKUPS=10
# cd "${BACKUP_DIR}"
# ls -t wiki_notes_backup_*.sql.gz | tail -n +$((MAX_BACKUPS + 1)) | xargs -r rm
# echo -e "${YELLOW}Cleaned up old backups (keeping last ${MAX_BACKUPS})${NC}"

echo -e "${GREEN}Backup location: ${BACKUP_DIR}/${BACKUP_FILE}${NC}"

