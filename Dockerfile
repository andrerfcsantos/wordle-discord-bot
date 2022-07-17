FROM golang:1.18.4

COPY . /app
WORKDIR /app

# Install app dependencies
RUN go get -v ./... && go build -v

# Run app
CMD ["go", "run", "main.go"]


