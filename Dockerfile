FROM golang:1.26.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /songju-agentic-dba .

EXPOSE 8080

CMD ["/songju-agentic-dba"]
