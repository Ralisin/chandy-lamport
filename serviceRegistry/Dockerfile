FROM golang:latest

WORKDIR /app

COPY ../go.mod go.sum ./
RUN go mod download


COPY ../config.json /app/config.json
COPY ../snapshotService /app/snapshotService

COPY serviceRegistry /app/serviceRegistry

CMD ["sh", "-c", "cd serviceRegistry && go run . -docker"]