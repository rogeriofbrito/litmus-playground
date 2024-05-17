#!/bin/sh

./litmusctl config set-account -e http://chaos-litmus-frontend-service.litmus.svc.cluster.local:9091 -u admin -p litmus --non-interactive=true
ADMIN_PROJECT_ID=$(./litmusctl get projects -o json | jq -r '.[] | select(.Name=="admin-project") | .ProjectID')
#TODO: validate if environment is already created
./litmusctl create chaos-environment --project-id="${ADMIN_PROJECT_ID}" --name="preproduction"
