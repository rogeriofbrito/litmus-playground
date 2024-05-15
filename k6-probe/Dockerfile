FROM grafana/k6
WORKDIR /k6-loadtest
COPY ./script.sh ./script.sh
# TODO: change to get script.js from git
COPY ./probe ./probe

# TODO: change to multistage build
USER root

RUN chmod +x script.sh

ENTRYPOINT [ "./script.sh" ]

# apk add kubernetes

# kubectl apply -f - <<EOF
# apiVersion: rbac.authorization.k8s.io/v1
# kind: Role
# metadata:
#   namespace: order-api-app
#   name: admin
# rules:
# - apiGroups: ["*"]
#   resources: ["*"]
#   verbs: ["*"]
# EOF

# kubectl create rolebinding default-admin --role=admin --serviceaccount=order-api-app:default --namespace=order-api-app