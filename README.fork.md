The Geth fork

The fork of an official golang implementation of the Ethereum protocol that saving blockchain data to the external database


## Installation

## Install system packages

```sh
sudo apt-get update
sudo apt-get install build-essential
sudo apt-get install git curl
sudo apt-get install postgresql postgresql-contrib
```

## Install Go

```sh
curl -O https://storage.googleapis.com/golang/go1.9.4.linux-amd64.tar.gz
tar -xvf go1.9.4.linux-amd64.tar.gz
sudo mv go /usr/local
nano ~/.profile
```

At the end of the file, add this line:

export PATH=$PATH:/usr/local/go/bin

```sh
source ~/.profile
```

Install Goose (is a database migration tool).

```sh
go get -u github.com/pressly/goose
```

## Set up database

```
sudo -u postgres psql
postgres=# CREATE DATABASE jsearch;
postgres=# CREATE USER jsearch WITH PASSWORD 'password';
postgres=# ALTER ROLE jsearch SET client_encoding TO 'utf8';
postgres=# ALTER ROLE jsearch SET default_transaction_isolation TO 'read committed';
postgres=# ALTER ROLE jsearch SET timezone TO 'UTC';
postgres=# GRANT ALL PRIVILEGES ON DATABASE jsearch TO jsearch;
postgres=# \q
```

## Init database

```sh
goose -dir  ~/go-ethereum/extdb/schema_migrations/ postgres "user=DB_USER_NAME dbname=DB_NAME sslmode=disable" up
```

## Clone project and create workdir

```
cd ~/
git clone https://github.com/jibrelnetwork/go-ethereum.git
git checkout feature/external_postgres_db
```

## Building the source

```
make geth
```

## Running geth

```
~/go-ethereum/build/bin/geth --syncmode full --extdb postgres://jsearch:password@localhost:5432/jsearch --cache 4096
```
