FROM golang:alpine

RUN mkdir /app
ADD . /app
WORKDIR /app
COPY ./transport/server .
COPY ./go.mod .
COPY ./go.sum .
COPY ./proto .

RUN go mod download && go mod verify

COPY . .

RUN go build
EXPOSE 50051

CMD ["go", "run", "."]