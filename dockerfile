FROM golang:1.23

WORKDIR /go/src

COPY ./src /go/src

# RUN go mod init charts
# RUN go get -u gorm.io/gorm
# RUN go get -u gorm.io/driver/mysql
# RUN go get -u github.com/labstack/echo/v4
# RUN go mod tidy

EXPOSE 1323