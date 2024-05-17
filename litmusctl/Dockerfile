FROM alpine:latest

WORKDIR /litmusctl

RUN apk add wget
RUN wget https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm64-1.6.0.tar.gz
RUN tar -xzvf litmusctl-linux-arm64-1.6.0.tar.gz

ENTRYPOINT [ "/litmusctl/litmusctl" ]
