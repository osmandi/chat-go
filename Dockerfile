FROM golang
//USER root
//RUN apt-get update && apt-get install git
EXPOSE 8080
WORKDIR go/src
COPY ./ .
CMD go get github.com/gorilla/websocket && cd src && go run main.go
