#!/bin/sh

k6 run ./probe/smoke-test-order-api/script.js -i 10
