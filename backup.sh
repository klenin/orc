#!/bin/bash

fileName="$(date +"%m_%d_%Y___%H_%M_%S").sql"

echo "$fileName"

pg_dump --host=localhost --port=5432 --username=admin --dbname=orc > "$fileName"
