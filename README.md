# litmus-playground

## Running minikube

1. Delete any existing minikube cluster

```bash
minikube delete
```

2. Start a new minikube cluster with `--insecure-registry` flag

```bash
minikube start --insecure-registry "10.0.0.0/24" --cpus 4
```

3. Enable Minikube Addons

```bash
minikube addons enable registry
minikube addons enable metrics-server
```

4. Run socat command to registry port

```bash
docker run --rm -it --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:$(minikube ip):5000"
```

5. Run port-forward to registry service

```bash
kubectl port-forward --namespace kube-system service/registry 5000:80
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

4. Run port-forward to Litmus frontend service

```bash
kubectl port-forward --namespace litmus service/chaos-litmus-frontend-service 8185:9091
```

## Install Chaos Infrastructure

1. Access Litmus ChaosCenter at [http://localhost:8185](http://localhost:8185) with username `admin` and password `litmus`

2. Create a Chaos Environment with name `env`

3. Create a Chaos Environment with name `infrastructure`

## Install order-api database

```bash
make order-api-k8s-postgres
```

## Run order-api database port-forward

```bash
kubectl port-forward --namespace postgres service/postgres-order-api 5432:5432
```

## Install order-api migration

```bash
make order-api-liquibase
```

## Install order-api app

```bash
make order-api-app
```

## Run order-api database port-forward

```bash
kubectl port-forward --namespace app service/order-api-service 8080:80
```

## Install k6

```bash
make k6
```

## Run k6 port-forward

```bash
kubectl port-forward --namespace k6 service/k6-service 5665:5665
```
