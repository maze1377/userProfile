FROM golang:1.16-bullseye as builder

ENV GO111MODULE=on
ENV GOPROXY=direct

ADD . /opt/

WORKDIR /opt

ADD go.mod .
ADD go.sum .
ADD Makefile .
RUN make dependencies

ADD . .
RUN make generate

FROM ubuntu:20.04 as runner

ENV DEBIAN_FRONTEND=noninteractive
ENV LANG=C.UTF-8 LC_ALL=C.UTF-8

RUN apt-get update --fix-missing && \
    apt-get upgrade -y && \
    apt-get install -y wget && \
    apt-get install -y ca-certificates && \
    apt-get install -y tzdata && \
    ln -sf /usr/share/zoneinfo/UTC /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata && \
    apt-get clean

COPY --from=builder /opt/userProfile /bin/userProfile

# Liveness
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.6 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

EXPOSE 10000
ENTRYPOINT ["/bin/userProfile"]
CMD ["version"]
