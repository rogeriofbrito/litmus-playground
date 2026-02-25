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
# TODO: add --version
--n cnpg-system \
--create-namespace \
cnpg/cloudnative-pg
```

### Install cnpg helm chart

```bash
source .env

helm upgrade --install cnpg ./helm/cnpg \
# TODO: add --version
-n $ORDER_CLUSTER_NS \
--create-namespace \
--set users.app.secrets.username=$ORDER_APP_USERNAME \
--set users.app.secrets.password=$ORDER_APP_PASSWORD \
--set users.root.secrets.username=$ORDER_ROOT_USERNAME \
--set users.root.secrets.password=$ORDER_ROOT_PASSWORD \
--set cluster.name=$ORDER_CLUSTER_NAME \
--set cluster.bootstrap.initdb.database=$ORDER_DATABASE_NAME \
--set cluster.bootstrap.initdb.owner=$ORDER_APP_USERNAME
```

### Install order-api helm chart

```bash
source .env

helm upgrade --install order-api ./helm/order-api \
-n order-api \
# TODO: add --version
--create-namespace \
--values ./helm-values/order-api.yaml \
--set secrets.DATABASE_HOST="$(echo -n "$ORDER_CLUSTER_NAME-rw.$ORDER_CLUSTER_NS.svc.cluster.local" | base64)" \
--set secrets.DATABASE_PORT="$(echo -n "5432" | base64)" \
--set secrets.DATABASE_NAME="$(echo -n $ORDER_DATABASE_NAME | base64)" \
--set secrets.DATABASE_USER="$(echo -n $ORDER_APP_USERNAME | base64)" \
--set secrets.DATABASE_PASSWORD="$(echo -n $ORDER_APP_PASSWORD | base64)"
```

### Install Prometheus helm chart

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm upgrade --install prometheus prometheus-community/prometheus \
# TODO: add --version
-n prometheus \
--create-namespace
```

### Install Grafana helm chart

```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm upgrade --install grafana grafana/grafana \
--version 10.5.15 \
-n grafana \
--create-namespace \
--values ./helm-values/grafana.yaml \
--set-file dashboards.default.golang-monitoring-dashboard.json=./grafana/dashboards/echo-framework-processes.json
```

## Obs

* When port-forwarding connection to a pg cluster pod, the connection is lost when testing connetion in Dbeaver. To use other functions it works well.

## References

* https://github.com/cloudnative-pg/charts
* https://cloudnative-pg.io/docs/devel/quickstart
* https://cloudnative-pg.io/docs/devel/installation_upgrade
