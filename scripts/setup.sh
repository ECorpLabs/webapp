sleep 30
sudo apt-get update -y
sudo apt install postgresql postgresql-contrib -y
sudo systemctl start postgresql.service

wget https://dl.google.com/go/go1.21.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xvf go1.21.1.linux-amd64.tar.gz 

echo export PATH=$PATH:/usr/local/go/bin >> ~/.bashrc
source ~/.bashrc


mv /tmp/webapp /home/admin/webapp
mv /tmp/users.csv /home/opt/users.csv

echo 'DB_HOST:"localhost"' >> .env 
echo 'DB_PORT:"5432"' >> .env 
echo 'DB_USER:"postgres"' >> .env 
echo 'DB_PASSWORD:"postgres"' >> .env 
echo 'DB_NAME:"dushyant"' >> .env 
echo 'APP_PORT:"8080"' >> .env 
echo 'FILE_PATH:"$HOME/opt/users.csv"' >> .env
