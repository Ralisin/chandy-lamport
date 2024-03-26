FROM golang:latest

WORKDIR /app

COPY ../go.mod go.sum ./
RUN go mod download

COPY ../config.json /app/config.json
COPY ../snapshotService /app/snapshotService

COPY peer /app/peer

CMD ["sh", "-c", "cd peer && go run . -docker"]