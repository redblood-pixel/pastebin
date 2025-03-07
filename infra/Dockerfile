FROM FROM golang:1.21-alpine 

RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash
RUN apt-get update
RUN apt-get install -y migrate