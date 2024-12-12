#!/bin/bash

export MON_CAVEAU_DB_USER=root
export MON_CAVEAU_DB_PASSWORD=mysql
export MON_CAVEAU_DB_HOST=localhost
export MON_CAVEAU_DB_PORT=3306
export MON_CAVEAU_DB_NAME=moncaveau

export DEBUG_MODE=true

cd ../src/

go run .