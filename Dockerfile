# use official Golang image
FROM golang:1.23

# set working directory
WORKDIR /app

# Copy the source code
COPY . .

RUN apk add --no-cache curl \
 && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz \
 && mv migrate /usr/bin/migrate \
 && chmod +x /usr/bin/migrate

# Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o api .

#EXPOSE the port
EXPOSE 8000

# Run the executable
CMD ["./api"]