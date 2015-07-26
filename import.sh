#!/bin/bash

echo Before importing close all processes that involve database orc.

echo Clearing the database to import data from a file.

psql --host=localhost --port=5432 --username=postgres -c "DROP DATABASE IF EXISTS orc;"

psql --host=localhost --port=5432 --username=postgres -c "CREATE DATABASE orc;"

psql --host=localhost --port=5432 --username=postgres -c "GRANT ALL PRIVILEGES ON DATABASE orc to admin;"

read -p "Enter the path to the file for importing: " pathToDatabase

psql -U postgres -d orc -f "$pathToDatabase"
