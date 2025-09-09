#!/bin/bash

# Скрипт для установки Liquibase и PostgreSQL драйвера

mkdir -p lib
if [ ! -f "liquibase-4.24.1.tar.gz" ]; then
    wget https://github.com/liquibase/liquibase/releases/download/v4.24.1/liquibase-4.24.1.tar.gz
fi

tar -xzf liquibase-4.24.1.tar.gz

sudo mv liquibase /usr/local/bin/
sudo chmod +x /usr/local/bin/liquibase

wget -O lib/postgresql-42.6.0.jar https://jdbc.postgresql.org/download/postgresql-42.6.0.jar

rm -f liquibase-4.24.1.tar.gz

liquibase --version
