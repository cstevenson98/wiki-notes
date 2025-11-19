#!/bin/bash

# Check backup health and status
# Can be used for monitoring or alerting

set -e

BACKUP_DIR="./backups"
MAX_AGE_DAYS=1

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "=== Wiki Notes Backup Health Check ==="
echo ""

# Check if backup directory exists
if [ ! -d "$BACKUP_DIR" ]; then
    echo -e "${RED}❌ Backup directory not found: $BACKUP_DIR${NC}"
    exit 1
fi

# Find latest backup
LATEST_BACKUP=$(ls -t ${BACKUP_DIR}/wiki_notes_backup_*.sql.gz 2>/dev/null | head -1)

if [ -z "$LATEST_BACKUP" ]; then
    echo -e "${RED}❌ No backups found in $BACKUP_DIR${NC}"
    echo ""
    echo "Create your first backup with:"
    echo "  ./scripts/backup-db.sh"
    exit 1
fi

# Calculate backup age
BACKUP_TIMESTAMP=$(stat -c %Y "$LATEST_BACKUP" 2>/dev/null || stat -f %m "$LATEST_BACKUP")
CURRENT_TIMESTAMP=$(date +%s)
BACKUP_AGE_SECONDS=$((CURRENT_TIMESTAMP - BACKUP_TIMESTAMP))
BACKUP_AGE_DAYS=$((BACKUP_AGE_SECONDS / 86400))
BACKUP_AGE_HOURS=$(((BACKUP_AGE_SECONDS % 86400) / 3600))

# Get backup size
BACKUP_SIZE=$(du -h "$LATEST_BACKUP" | cut -f1)

# Count total backups
TOTAL_BACKUPS=$(ls -1 ${BACKUP_DIR}/wiki_notes_backup_*.sql.gz 2>/dev/null | wc -l)

# Check backup age
if [ $BACKUP_AGE_DAYS -gt $MAX_AGE_DAYS ]; then
    echo -e "${YELLOW}⚠️  WARNING: Latest backup is ${BACKUP_AGE_DAYS} days old${NC}"
    STATUS="WARNING"
    EXIT_CODE=1
elif [ $BACKUP_AGE_DAYS -eq 0 ] && [ $BACKUP_AGE_HOURS -lt 1 ]; then
    echo -e "${GREEN}✅ EXCELLENT: Very recent backup${NC}"
    STATUS="EXCELLENT"
    EXIT_CODE=0
else
    echo -e "${GREEN}✅ OK: Backup is recent${NC}"
    STATUS="OK"
    EXIT_CODE=0
fi

echo ""
echo "Latest backup:"
echo "  File: $(basename $LATEST_BACKUP)"
echo "  Size: $BACKUP_SIZE"
if [ $BACKUP_AGE_DAYS -eq 0 ]; then
    echo "  Age: ${BACKUP_AGE_HOURS} hours"
else
    echo "  Age: ${BACKUP_AGE_DAYS} days, ${BACKUP_AGE_HOURS} hours"
fi
echo "  Total backups: $TOTAL_BACKUPS"

# Check disk space
BACKUP_DIR_SIZE=$(du -sh "$BACKUP_DIR" | cut -f1)
echo ""
echo "Backup directory size: $BACKUP_DIR_SIZE"

# Check if database is running
if docker ps --format '{{.Names}}' | grep -q "wiki-notes-db-1"; then
    echo -e "${GREEN}Database container: Running ✓${NC}"
else
    echo -e "${YELLOW}Database container: Not running${NC}"
fi

echo ""
echo "=== Status: $STATUS ==="

exit $EXIT_CODE

