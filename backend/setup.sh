#!/bin/bash

# Script for setting up the backend side of our application

echo "Beginning setup for PostgreScrutiniser...\n"
echo "Note that port 9090 has to be open for the application to be accessible"

# Create necessasry folders and setup our main application user
mkdir -p /usr/local/postgrescrutiniser/backups/
mkdir -p /usr/local/postgrescrutiniser/confs/

APPUSER="postgrescrutiniser"
HOMEDIR="/usr/local/$APPUSER/"
SCRUTINISER_PASSWORD=`date +%s | sha256sum | base64 | head -c 16`
echo $SCRUTINISER_PASSWORD | passwd --stdin $APPUSER
useradd -m -d $HOMEDIR $APPUSER

echo "postgrescrutiniser  ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

echo "Main application user $APPUSER created with password: $SCRUTINISER_PASSWORD"
echo "These credentials should be used to connect to our application"

# Save postgres database's main user's credentials into a .pgpass file
echo "What is the name of the main PostgreSql database user?"
read POSTGREUSER
echo "What is its password?"
read PASSWORD

PORT=$(lsof -i -P -n | awk '/post/ {print $9}' | awk -F ':' '{print $2}' | head -1)
if [ $PORT == '' ]
then
  echo "What port is PostgreSql running on?"
  read PORT
fi

echo "*:$PORT:*:$POSTGREUSER:$PASSWORD" > $HOMEDIR/.pgpass
chown -R $APPUSER. $HOMEDIR
chmod 0600 $HOMEDIR/.pgpass

echo "Main PostgreSql database's user credentials saved in: $HOMEDIR/.pgpass"

# Set up backend executable
echo "Setting up PostgreScrutiniser executable..."
go build -o postgrescrutiniser
cp -p postgrescrutiniser /usr/local/postgrescrutiniser/
cp -p postgrescrutiniser.service /etc/systemd/system

SECRET=$(echo $RANDOM | md5sum | head -c 20)
echo "BACKEND_PORT=9090" > dev.env
echo "JWT_SECRET_KEY=${SECRET}" >> dev.env
mv dev.env /usr/local/postgrescrutiniser/

echo "Setup completed without any errors."