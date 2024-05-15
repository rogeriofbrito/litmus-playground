FROM grafana/k6
WORKDIR /k6-loadtest
COPY ./script.sh ./script.sh
# TODO: change to get script.js from git
COPY ./probe ./probe

# TODO: change to multistage build
USER root

RUN chmod +x script.sh

ENTRYPOINT [ "./script.sh" ]
