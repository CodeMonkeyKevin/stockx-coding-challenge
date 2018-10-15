#!/bin/sh

set -e

echo "Starting the Postgresql daemon"
service postgresql start

export APP_DB_NAME=stockxcc_kp
export APP_DB_USERNAME=docker
export APP_DB_PASSWORD=docker

echo "Starting Shoe API"
/bin/stockx
