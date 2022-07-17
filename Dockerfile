FROM golang:1.18.4

COPY . /app
WORKDIR /app

# Install app dependencies and build the app
RUN go mod tidy && go build -o wordle-discord-bot

# Run app
CMD ["./wordle-discord-bot"]


