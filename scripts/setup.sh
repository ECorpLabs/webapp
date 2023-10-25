#!/bin/bash

sleep 30
sudo apt-get update -y
sudo apt-get upgrade -y
# sudo apt install postgresql postgresql-contrib -y
# sudo systemctl start postgresql.service

wget https://dl.google.com/go/go1.21.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xvf go1.21.1.linux-amd64.tar.gz 

echo export PATH=$PATH:/usr/local/go/bin >> ~/.bashrc
source ~/.bashrc

# sudo -u postgres psql -c "ALTER USER postgres WITH PASSWORD 'postgres';"

sudo groupadd csye6225
sudo useradd -s /bin/false -g csye6225 -d /opt/csye6225 -m csye6225
# sudo cp csye6225.service /etc/systemd/system
systemctl daemon-reload
sudo systemctl enable csye6225
sudo systemctl start csye6225
sudo systemctl restart csye6225
sudo systemctl stop csye6225

mv /tmp/webapp /opt/webapp
sudo mv /tmp/users.csv /opt/users.csv

# {
#     echo 'DB_HOST:"localhost"'
#     echo 'DB_PORT:"5432"'
#     echo 'DB_USER:"postgres"'
#     echo 'DB_PASSWORD:"postgres"'
#     echo 'DB_NAME:"dushyant"'
#     echo 'APP_PORT:"8080"'
#     echo 'FILE_PATH:"/opt/users.csv"'
# } >> .env