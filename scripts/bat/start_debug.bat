@echo off

set MON_CAVEAU_DB_USER=user
set MON_CAVEAU_DB_PASSWORD=password
set MON_CAVEAU_DB_HOST=host
set MON_CAVEAU_DB_PORT=port
set MON_CAVEAU_DB_NAME=name

set USE_TLS=false
set CERT_FILE=""
set KEY_FILE=""
set DOMAIN_NAME=""

set DEBUG_MODE=true
set USE_FILESYSTEM_FRONTEND=true

set ACCOUNT_KEY_SECRET="ThisIsForDebugPurposes"

set ACTIVITY_FLUSH_INTERVAL=120000

cd ../src/

go run .