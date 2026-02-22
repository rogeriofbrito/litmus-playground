# litmus-playgraound

## Create kind cluster (for local deploy)

Create a kind cluster to deploy all services:

```bash
kind create cluster
```

## Setup environment variables

Create a `.env` file at root directory using `.env.template` as template.

## Helm Charts

### Install cloud-native-pg operator

```bash
helm repo add cnpg https://cloudnative-pg.github.io/charts
helm upgrade --install cnpg-system \
  --namespace cnpg-system \
  --create-namespace \
  cnpg/cloudnative-pg
```

### Install cnpg helm chart

```bash
source .env

helm upgrade --install cnpg ./helm/cnpg \
-n cnpg \
--create-namespace \
--set users.app.secrets.username=$APP_USER_USERNAME \
--set users.app.secrets.password=$APP_USER_PASSWORD \
--set users.root.secrets.username=$SUPER_USER_USERNAME \
--set users.root.secrets.password=$SUPER_USER_PASSWORD \
--set cluster.bootstrap.initdb.database=$APP_USER_DATABASE \
--set cluster.bootstrap.initdb.owner=$APP_USER_USERNAME
```

### Install order-api helm chart

```bash
source .env

helm upgrade --install order-api ./helm/order-api \
-n order-api \
--create-namespace \
--set secrets.DATABASE_HOST="$(echo -n "pg-cluster-rw.cnpg.svc.cluster.local" | base64)" \
--set secrets.DATABASE_PORT="$(echo -n "5432" | base64)" \
--set secrets.DATABASE_NAME="$(echo -n $APP_USER_DATABASE | base64)" \
--set secrets.DATABASE_USER="$(echo -n $APP_USER_USERNAME | base64)" \
--set secrets.DATABASE_PASSWORD="$(echo -n $APP_USER_PASSWORD | base64)"
```

## Obs

* When port-forwarding connection to a pg cluster pod, the connection is lost when testing connetion in Dbeaver. To use other functions it works well.

## References

* https://github.com/cloudnative-pg/charts
* https://cloudnative-pg.io/docs/devel/quickstart
* https://cloudnative-pg.io/docs/devel/installation_upgrade