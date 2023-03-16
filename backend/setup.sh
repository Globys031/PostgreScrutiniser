#!/bin/bash

# For now I'll just have it create postgrescrutiniser, 
# setup its home directory and create .pgpass file.
# Later it should also setup the entirety of our application

echo "Beginning setup for PostgreScrutiniser.\n"

APPUSER="postgrescrutiniser"
HOMEDIR="/usr/local/$APPUSER/"
useradd -m -d $HOMEDIR $APPUSER

echo "What is the name of the main PostgreSql database user?"
read POSTGREUSER
echo "What is its password?"
read PASSWORD

echo "postgrescrutiniser  ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

PORT=$(lsof -i -P -n | awk '/post/ {print $9}' | awk -F ':' '{print $2}' | head -1)
if [ $PORT == '' ]
then
  echo "What port is PostgreSql running on?"
  read PORT
fi

echo "*:$PORT:*:$POSTGREUSER:$PASSWORD" > $HOMEDIR/.pgpass
chown $APPUSER $HOMEDIR/.pgpass
chmod 0600 $HOMEDIR/.pgpass