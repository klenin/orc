@echo off

echo "Before importing close all processes that involve database orc"

psql --host=localhost --port=5432 --username=postgres -c "DROP DATABASE IF EXISTS orc;"

psql --host=localhost --port=5432 --username=postgres -c "CREATE DATABASE orc;"

psql --host=localhost --port=5432 --username=postgres -c "GRANT ALL PRIVILEGES ON DATABASE orc to admin;"

echo "Enter password of remote database: "

pg_dump --host=<host> --port=<port> --username=<username> --dbname=<dbname> > output.sql

cat output.sql | sed 's/<username>/admin/g' > dump.sql

psql -U postgres -d orc -f dump.sql

psql -U postgres -d orc
