@echo off

psql --host=localhost --port=5432 --username=postgres -c "DROP DATABASE IF EXISTS orc;"

psql --host=localhost --port=5432 --username=postgres -c "CREATE DATABASE orc;"

psql --host=localhost --port=5432 --username=postgres -c "GRANT ALL PRIVILEGES ON DATABASE orc to admin;"

echo "Enter password of remote database: "

pg_dump --host=ec2-23-23-225-50.compute-1.amazonaws.com --port=5432 --username=hduspsokjkhsmj --dbname=d8lt1tbga2v27l > output.sql

cat output.sql | sed 's/hduspsokjkhsmj/admin/g' > new.sql

psql -U postgres -d orc -f new.sql

psql -U postgres -d orc
