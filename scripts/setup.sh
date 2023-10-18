#!/bin/bash

sleep 30
sudo apt-get update -y
sudo apt-get upgrade -y
sudo apt install postgresql postgresql-contrib -y
sudo systemctl start postgresql.service

wget https://dl.google.com/go/go1.21.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xvf go1.21.1.linux-amd64.tar.gz 

echo export PATH=$PATH:/usr/local/go/bin >> ~/.bashrc
source ~/.bashrc


mv /tmp/webapp /home/admin/webapp
mv /tmp/users.csv /opt/users.csv

{
    echo 'DB_HOST:"localhost"'
    echo 'DB_PORT:"5432"'
    echo 'DB_USER:"postgres"'
    echo 'DB_PASSWORD:"postgres"'
    echo 'DB_NAME:"dushyant"'
    echo 'APP_PORT:"8080"'
    echo 'FILE_PATH:"/opt/users.csv"'
} >> .env