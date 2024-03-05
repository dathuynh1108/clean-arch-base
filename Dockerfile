FROM golang:alpine as builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux

WORKDIR /service/app/

COPY . .

RUN go mod download
RUN go build -a -installsuffix cgo -o ./main .

FROM alpine

WORKDIR /service/app/

RUN  apk add --no-cache bash \
        \
        && apk add --no-cache tzdata \
        \
        && cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime \
        \
        && echo "Asia/Ho_Chi_Minh" > /etc/timezone

COPY --from=builder /service/app/main .

ENTRYPOINT ["/service/app/main"]
CMD ["--config", "/service/app/configs/config.yaml"]