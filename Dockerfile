FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY *.go ./
RUN go build -o /simple-backend-server

ENV LISTEN_ON=$LISTEN_ON
ENV C_BACKEND_URL=$C_BACKEND_URL
ENV A_BACKEND_URL=$A_BACKEND_URL

EXPOSE 8080

CMD [ "/simple-backend-server" ]