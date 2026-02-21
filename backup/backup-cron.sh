#!/bin/sh

DATABASE_PATH="/data/equipment-watchdog.db"
BACKUP_BASE_DIR="/backups"

YEAR_MONTH=$(date +"%Y%m")
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

BACKUP_DIR="${BACKUP_BASE_DIR}/${YEAR_MONTH}"
BACKUP_FILE="equipment-watchdog_backup_${TIMESTAMP}.sql"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

echo "$(date '+%Y-%m-%d %H:%M:%S') - Starting backup..."

# Check if database exists
if [ ! -f "$DATABASE_PATH" ]; then
    echo "$(date '+%Y-%m-%d %H:%M:%S') - ERROR: Database not found at $DATABASE_PATH"
    exit 1
fi

# Create backup using sqlite3 dump
if sqlite3 "$DATABASE_PATH" ".dump" > "${BACKUP_DIR}/${BACKUP_FILE}"; then
    echo "$(date '+%Y-%m-%d %H:%M:%S') - Backup created successfully: ${YEAR_MONTH}/${BACKUP_FILE}"
else
    echo "$(date '+%Y-%m-%d %H:%M:%S') - ERROR: Backup failed"
    exit 1
fi
