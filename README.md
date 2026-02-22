# litmus-playgraound

## Create kind cluster (for local deploy)

Create a kind cluster to deploy all services:

```bash
kind create cluster
```

## Setup environment variables

Create a `.env` file at root directory using `.env.template` as template.

## cloud-native-pg

### Install cloud-native-pg operator

```bash
helm repo add cnpg https://cloudnative-pg.github.io/charts
helm upgrade --install cnpg \
  --namespace cnpg-system \
  --create-namespace \
  cnpg/cloudnative-pg
```

### Create Cluster

```bash
source .env

kubectl apply -f ./k8s/pg-cluster/1-namespace.yaml
envsubst < ./k8s/pg-cluster/2-app-user-secret.yaml | kubectl apply -f -
envsubst < ./k8s/pg-cluster/3-super-user-secret.yaml | kubectl apply -f -
envsubst < ./k8s/pg-cluster/4-pg-cluster.yaml | kubectl apply -f -

kubectl get -n pg cluster
```

### Obs

* When port-forwarding connection to a pg cluster pod, the connection is lost when testing connetion in Dbeaver. To use other functions it works well.

### References

https://github.com/cloudnative-pg/charts
https://cloudnative-pg.io/docs/devel/quickstart
https://cloudnative-pg.io/docs/devel/installation_upgrade

## kubernetes-playground

### Install kubernetes-playground helm chart

```bash
source .env

helm upgrade --install order-api ./helm/order-api \
-n order-api \
--create-namespace \
--set secrets.DATABASE_HOST="$(echo -n "pg-cluster-rw.pg.svc.cluster.local" | base64)" \
--set secrets.DATABASE_PORT="$(echo -n "5432" | base64)" \
--set secrets.DATABASE_NAME="$(echo -n $APP_USER_DATABASE | base64)" \
--set secrets.DATABASE_USER="$(echo -n $APP_USER_USERNAME | base64)" \
--set secrets.DATABASE_PASSWORD="$(echo -n $APP_USER_PASSWORD | base64)"
```
