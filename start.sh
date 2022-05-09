#!/bin/sh
set -e

 echo "Running db migrations"
 echo "This is DB_URL $DB_URL"
 /app/migrate -path /app/migration -database $DB_URL -verbose up
 
 echo "Starting search service . . ."
 ./main