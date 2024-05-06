REGISTRY := localhost:5000
LIQUIBASE_ORDER_API_IMAGE_NAME := order-api-liquibase
LIQUIBASE_ORDER_API_TAG := latest
ORDER_API_IMAGE_NAME := order-api
ORDER_API_TAG := latest
K6_IMAGE_NAME := k6
K6_TAG := latest

# order-api K8S Postgres

order-api-k8s-postgres:
	kubectl apply -f ./order-api/k8s/postgres/namespace.yaml
	kubectl apply -f ./order-api/k8s/postgres/pvc.yaml
	kubectl apply -f ./order-api/k8s/postgres/deployment.yaml
	kubectl apply -f ./order-api/k8s/postgres/service.yaml

# order-api K8S Liquibase

order-api-liquibase: order-api-docker-liquibase order-api-k8s-liquibase

order-api-docker-liquibase:
	docker build --no-cache -t $(LIQUIBASE_ORDER_API_IMAGE_NAME) ./order-api/migration
	docker tag $(LIQUIBASE_ORDER_API_IMAGE_NAME) $(REGISTRY)/$(LIQUIBASE_ORDER_API_IMAGE_NAME):$(LIQUIBASE_ORDER_API_TAG)
	docker push $(REGISTRY)/$(LIQUIBASE_ORDER_API_IMAGE_NAME):$(LIQUIBASE_ORDER_API_TAG)

order-api-k8s-liquibase:
	REGISTRY=$(REGISTRY) IMAGE=$(LIQUIBASE_ORDER_API_IMAGE_NAME) TAG=$(LIQUIBASE_ORDER_API_TAG) envsubst < ./order-api/k8s/migration/pod.yaml | kubectl delete --force -f - || true
	REGISTRY=$(REGISTRY) IMAGE=$(LIQUIBASE_ORDER_API_IMAGE_NAME) TAG=$(LIQUIBASE_ORDER_API_TAG) envsubst < ./order-api/k8s/migration/pod.yaml | kubectl apply --force -f -

# order-api K8S App

order-api-app: order-api-docker-app order-api-k8s-app

order-api-docker-app:
	docker build --no-cache -t $(ORDER_API_IMAGE_NAME) ./order-api
	docker tag $(ORDER_API_IMAGE_NAME) $(REGISTRY)/$(ORDER_API_IMAGE_NAME):$(ORDER_API_TAG)
	docker push $(REGISTRY)/$(ORDER_API_IMAGE_NAME):$(ORDER_API_TAG)

order-api-k8s-app:
	kubectl apply --force -f ./order-api/k8s/namespace.yaml
	REGISTRY=$(REGISTRY) IMAGE=$(ORDER_API_IMAGE_NAME) TAG=$(ORDER_API_TAG) envsubst < ./order-api/k8s/deployment.yaml | kubectl delete --force -f - || true
	REGISTRY=$(REGISTRY) IMAGE=$(ORDER_API_IMAGE_NAME) TAG=$(ORDER_API_TAG) envsubst < ./order-api/k8s/deployment.yaml | kubectl apply --force -f -
	kubectl apply --force -f ./order-api/k8s/service.yaml

# K6

k6: k6-docker k6-k8s

k6-docker:
	docker build --no-cache -t $(K6_IMAGE_NAME) ./k6-loadtest
	docker tag $(K6_IMAGE_NAME) $(REGISTRY)/$(K6_IMAGE_NAME):$(K6_TAG)
	docker push $(REGISTRY)/$(K6_IMAGE_NAME):$(K6_TAG)

k6-k8s:
	kubectl apply --force -f ./order-api/k8s/k6/namespace.yaml
	REGISTRY=$(REGISTRY) IMAGE=$(K6_IMAGE_NAME) TAG=$(K6_TAG) envsubst < ./order-api/k8s/k6/deployment.yaml | kubectl delete --force -f - || true
	REGISTRY=$(REGISTRY) IMAGE=$(K6_IMAGE_NAME) TAG=$(K6_TAG) envsubst < ./order-api/k8s/k6/deployment.yaml | kubectl apply --force -f -
	kubectl apply --force -f ./order-api/k8s/k6/service.yaml
