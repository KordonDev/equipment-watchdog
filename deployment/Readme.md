### Create backup file
sqlite3 data/equipment-watchdog.db .dump >> backup.sql

### Import backup data
There are many ways to do this, one way is:

sqlite3 data/equipment-watchdog.db

Followed by:

sqlite> .read backup.sql
