# litmus-playground with Kind

## Running Kind

1. Create a new Kind cluster if it doesn't exists

```bash
kind create cluster --config kind-config.yaml
```

2. Set `kind-kind` as current kubectl context

```bash
kubectl config use-context kind-kind
```

3. Create metrics-server in kind cluster

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.5.0/components.yaml
kubectl patch deployment metrics-server -n kube-system --patch "$(cat metric-server-patch.yaml)"
```

## Install Litmus

1. Create a namespace in which Litmus will be installed

```bash
kubectl create ns litmus
```

2. Add Litmus Helm Chart

```bash
helm repo add litmuschaos https://litmuschaos.github.io/litmus-helm/
```

3. Install Litmus Helm Chart

If the command will be executed in a ARM-based machine, specify a Mongo DB compatible image:

```bash
helm install chaos litmuschaos/litmus --namespace=litmus \
--set portal.frontend.service.type=NodePort \
--set mongodb.image.registry=docker.io \
--set mongodb.image.repository=zcube/bitnami-compat-mongodb \
--set mongodb.image.tag=6.0.5 \
--set mongodb.volumePermissions.image.registry=docker.io \
--set mongodb.volumePermissions.image.repository=bitnami/os-shell \
--set mongodb.volumePermissions.image.tag=12-debian-12-r19
```

In case of others CPUs, you can use the default Mongo DB image:

```bash
helm install chaos litmuschaos/litmus --namespace=litmus \
--set portal.frontend.service.type=NodePort
```

4. Apply patch to change NodePort port

```bash
kubectl patch service chaos-litmus-frontend-service -n litmus --patch "$(cat litmus-frontend-service-patch.yaml)"
```

## Install Chaos Infrastructure

1. Access Litmus ChaosCenter at [http://localhost:8185](http://localhost:8185) with username `admin` and password `litmus`

2. Create a Pre-Production Chaos Environment with name `preproduction`

3. Install Litmus Agent Helm Chart (remember replace `<project-id>` value)

```bash
helm install litmus-agent litmuschaos/litmus-agent \
--namespace litmus \
--set "INFRA_NAME=helm-agent" \
--set "INFRA_DESCRIPTION=helm-agent" \
--set "LITMUS_URL=http://chaos-litmus-frontend-service.litmus.svc.cluster.local:9091" \
--set "LITMUS_BACKEND_URL=http://chaos-litmus-server-service.litmus.svc.cluster.local:9002" \
--set "LITMUS_USERNAME=admin" \
--set "LITMUS_PASSWORD=litmus" \
--set "LITMUS_PROJECT_ID=<project-id>" \
--set "LITMUS_ENVIRONMENT_ID=preproduction"
```

## Install order-api

```bash
make order-api
```

## Install probe-api

```bash
make probe-api
```

## Install k6-probe

```bash
make k6-probe
```

## Install k6-loadtest

```bash
make k6-loadtest
```