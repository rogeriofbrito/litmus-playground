#!/bin/sh

k6 run -e ORDER_HOST=order-api-app-service.order-api-app.svc.cluster.local -e ORDER_PORT=8080 ./probe/smoke-test-order-api/script.js -i 10
