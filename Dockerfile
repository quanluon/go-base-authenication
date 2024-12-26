FROM golang:1.21.3-alpine

WORKDIR /app

# Install make
# RUN apt update && apt install -y make

COPY . .

RUN go mod download

COPY . .

EXPOSE 8888

CMD ["go", "run", "main.go"]