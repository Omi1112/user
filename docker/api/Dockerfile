FROM golang:latest

WORKDIR /go/src/

RUN apt-get update \
  && apt-get install

RUN go get -u github.com/gin-gonic/gin \
  && go get github.com/jinzhu/gorm \
  && go get github.com/go-sql-driver/mysql \
  && go get github.com/codegangsta/gin \
  && apt-get install -y mariadb-client

CMD ["sh", "./start_app.sh"]
