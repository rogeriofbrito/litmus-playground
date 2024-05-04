# litmus-playground

## Running minikube with registry addon

1. Delete any existing minikube cluster

```bash
minikube delete
```

2. Start a new minikube cluster with `--insecure-registry` flag

```bash
minikube start --insecure-registry "10.0.0.0/24"
```

3. Enable the registry addon on Minikube cluster

```bash
minikube addons enable registry
```

4. Run socat command

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
--set mongodb.image.registry=ghcr.io/zcube \
--set mongodb.image.repository=bitnami-compat/mongodb \
--set mongodb.image.tag=6.0.5
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
