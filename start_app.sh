# !/bin/bash

# MySQLサーバーが起動するまで待機する
until mysqladmin ping -h mysql -P 3306 --silent; do
  echo 'waiting for mysqld to be connectable...'
  sleep 2
done

echo "app is starting...!"
exec gin -i --all run main.go
