#!/bin/bash

# Script for setting up th entirety of our application

echo "Beginning setup for PostgreScrutiniser.\n"

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
chown $APPUSER $HOMEDIR/.pgpass
chmod 0600 $HOMEDIR/.pgpass

echo "Main PostgreSql database's user credentials saved in: $HOMEDIR/.pgpass"