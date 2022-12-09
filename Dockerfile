FROM golang:1.19.3-buster as builder

RUN apt update && apt install curl unzip -y

ENV GOPATH=/go

WORKDIR /src

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download || true

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

#Build docker
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./out/server ./cmd/server


######## Start a new stage from scratch #######
FROM alpine:3.13

RUN apk --no-cache add ca-certificates tzdata htop tini bash curl

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /src/out/server /bin/

# Copy resources file
COPY --from=builder /src/resources /resources/

# Expose ports
EXPOSE 8088 8088

# Command to run the executable
CMD ["tini", "--", "/bin/server"]