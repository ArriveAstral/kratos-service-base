FROM golang:1.15 as builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 10001
EXPOSE 10002

VOLUME /applications/kratos-service-base/user/configs

CMD ["./user", "-conf", "/applications/kratos-service-base/user/configs/config.yaml"]
