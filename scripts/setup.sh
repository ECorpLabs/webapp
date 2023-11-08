#!/bin/bash

sleep 30
sudo apt-get update -y
sudo apt-get upgrade -y
# sudo apt install postgresql postgresql-contrib -y
# sudo systemctl start postgresql.service

wget https://amazoncloudwatch-agent.s3.amazonaws.com/debian/amd64/latest/amazon-cloudwatch-agent.deb
sudo dpkg -i -E ./amazon-cloudwatch-agent.deb

wget https://dl.google.com/go/go1.21.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xvf go1.21.1.linux-amd64.tar.gz 

echo export PATH=$PATH:/usr/local/go/bin >> ~/.bashrc
source ~/.bashrc

# sudo -u postgres psql -c "ALTER USER postgres WITH PASSWORD 'postgres';"

sudo groupadd csye6225
sudo useradd -s /bin/false -g csye6225 -d /opt/csye6225 -m csye6225

sudo mv /tmp/webapp /opt/csye6225/
sudo mv /tmp/users.csv /opt/
sudo mv /tmp/system.service /etc/systemd/system/
sudo mv /tmp/cloudwatch-config.json /opt/aws/amazon-cloudwatch-agent/etc/
sudo mkdir /var/log/webapp/
sudo touch /var/log/webapp/app.log
sudo chown csye6225:csye6225 /var/log/webapp/app.log

sudo touch /opt/csye6225/.env
sudo chown csye6225:csye6225 /opt/csye6225/.env

sudo systemctl daemon-reload
sudo systemctl enable system
# sudo systemctl start system
# sudo systemctl restart csye6225
# sudo systemctl stop csye6225




# {
#     echo 'DB_HOST:"localhost"'
#     echo 'DB_PORT:"5432"'
#     echo 'DB_USER:"postgres"'
#     echo 'DB_PASSWORD:"postgres"'
#     echo 'DB_NAME:"dushyant"'
#     echo 'APP_PORT:"8080"'
#     echo 'FILE_PATH:"/opt/users.csv"'
# } >> .env