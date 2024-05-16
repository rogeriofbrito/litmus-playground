#!/bin/sh

k6 run ./probe/$PROBE_ID/script.js --vus 1 --duration $DURATION
