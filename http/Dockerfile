FROM golang:alpine

RUN mkdir /app
ADD . /app
WORKDIR /app
COPY ./http .
COPY ./go.mod .
COPY ./go.sum .
COPY ./proto .
COPY ./assets .

RUN go mod download && go mod verify

COPY . .

RUN go build
EXPOSE 8080

CMD ["go", "run", "."]