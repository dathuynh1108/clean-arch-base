FROM golang:alpine as builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux

WORKDIR /service/app/

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .

RUN go build -a -installsuffix cgo -o ./main .

FROM alpine

WORKDIR /service/app/

RUN apk add --no-cache bash \
        \
        && apk add --no-cache tzdata \
        \
        && cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime \
        \
        && echo "Asia/Ho_Chi_Minh" > /etc/timezone

COPY --from=builder /service/app/main .
COPY --from=builder /service/app/configs/config.toml /service/app/configs/config.toml

ARG VERSION
RUN echo $VERSION > ./version.txt

ENTRYPOINT ["/service/app/main"]
EXPOSE 8080
CMD ["default", "-c", "/service/app/configs/config.toml"]
