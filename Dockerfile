FROM golang:1.18.3 AS development
WORKDIR /dd_project
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/cespare/reflex
EXPOSE 8080
CMD reflex -g '*.go' go run . --start-service