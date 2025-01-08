#!/bin/bash

export MON_CAVEAU_DB_USER=user
export MON_CAVEAU_DB_PASSWORD=password
export MON_CAVEAU_DB_HOST=host
export MON_CAVEAU_DB_PORT=port
export MON_CAVEAU_DB_NAME=name

export USE_TLS=false
export CERT_FILE=""
export KEY_FILE=""
export DOMAIN_NAME=domain.name

export DEBUG_MODE=true

export ACCOUNT_KEY_SECRET="ThisIsForDebugPurposes"

cd ../../src/

go run .