#!/bin/sh

k6 run ./probe/$PROBE_ID/script.js -i 10
