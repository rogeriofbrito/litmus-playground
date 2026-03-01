FROM grafana/k6:latest

COPY script.js /test/script.js

ENTRYPOINT ["k6"]
CMD ["run", "/test/script.js"]
