@echo off

set MON_CAVEAU_DB_USER=user
set MON_CAVEAU_DB_PASSWORD=password
set MON_CAVEAU_DB_HOST=host
set MON_CAVEAU_DB_PORT=port
set MON_CAVEAU_DB_NAME=name

set DEBUG_MODE=false

cd ../src/

go run .