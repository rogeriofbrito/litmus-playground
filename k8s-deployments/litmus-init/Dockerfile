FROM litmusctl

COPY script.sh script.sh

RUN chmod +x script.sh
RUN apk add jq

ENTRYPOINT [ "./script.sh" ]
