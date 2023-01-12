#!/bin/bash
# This script will reset the redis database
docker compose stop redis
docker compose rm -f redis
docker compose up -d redis