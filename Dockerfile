FROM golang:1.11-alpine AS builder

# Install some dependencies needed to build the project
RUN apk add bash git

#RUN apt-get update && apt-get install -y wget

# create a working directory
WORKDIR /neb/github.com/sovikc/neb

# Force the go compiler to use modules
ENV GO111MODULE=on

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download


# add source code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o neb .

FROM scratch
WORKDIR /neb
COPY --from=builder /neb/github.com/sovikc/neb/ .
EXPOSE 8080


ENTRYPOINT [ "./neb" ]