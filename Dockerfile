FROM golang:1.18

WORKDIR /data

COPY . .

RUN go mod vendor

RUN go build --mod=vendor --trimpath -o app ./init
ENV SCHEMA_PATH=./schema
ENV DB_HOST=localhost
ENV DB_PORT=5438
ENV DB_USERNAME=postgres
ENV DB_PASS=postgres
ENV DB_NAME=postgres
ENV API_PORT=8083

CMD ["./app"]