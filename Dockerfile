FROM golang:1.19

RUN apt update

WORKDIR /src

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go mod download

COPY . .

RUN go build -race

ENTRYPOINT ["./XLangConstants", "-input=/input", "-output=/output"]