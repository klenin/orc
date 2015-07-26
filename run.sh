#!/bin/bash

export DATABASE_URL="user=admin host=localhost dbname=orc password=admin sslmode=disable"

export PORT="6543"

read -p "Clear the database and run the system with test data [y/n]: " loadTestDataOrNot

if [[ "$loadTestDataOrNot" == "y" || "$loadTestDataOrNot" == "Y" ]];
then go build && orc.exe -test-data=true
else go build && orc.exe -test-data=false
fi
