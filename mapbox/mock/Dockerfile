FROM golang:1.17

WORKDIR /go/src/mapbox

ENV GOOS=linux
ENV GOARCH=amd64

COPY . .

RUN go build -o mapbox-${GOOS}-${GOARCH}

ENTRYPOINT ./mapbox-${GOOS}-${GOARCH}


