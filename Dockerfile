FROM golang:alpine AS builder
WORKDIR /dd_project
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main main.go
#RUN go install github.com/cespare/reflex
#CMD reflex -g '*.go' go run . --start-service


FROM alpine
WORKDIR /dd_project
COPY --from=builder /dd_project/main .
COPY . .
COPY .env .
COPY start.sh .
COPY wait-for.sh .


EXPOSE 8080
CMD [ "/dd_project/main" ]
ENTRYPOINT [ "/dd_project/start.sh" ]
