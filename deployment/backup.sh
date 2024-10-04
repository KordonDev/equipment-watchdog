#!/bin/bash

DATABASE_PATH="data/equipment-watchdog.db"
BACKUP_DESTINATION="."

CURRENT_FOLDER_NAME=$(basename "$PWD")

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

BACKUP_FILE="${CURRENT_FOLDER_NAME}_backup_${TIMESTAMP}.db"

sqlite3 "$DATABASE_PATH" ".backup '$BACKUP_FILE'"

mv "$BACKUP_FILE" "$BACKUP_DESTINATION"

echo "created backup: $BACKUP_FILE"

