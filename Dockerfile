FROM golang:1.17

WORKDIR /go/src/sport

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=1

ENV SPORT_SQLITE_PATH=/tmp/sport/sport.sqlite
ENV SPORT_SESSION_STORE_PATH=/tmp/sport/sessions
ENV SPORT_UPLOAD_FOLDER=/tmp/sport/uploads

RUN mkdir -p ${SPORT_SESSION_STORE_PATH} ${SPORT_UPLOAD_FOLDER}
RUN apt-get update && apt-get install -y sqlite3

COPY . .

RUN go build -o sport
RUN cp ./sport /usr/local/bin/sport

CMD /usr/local/bin/sport
