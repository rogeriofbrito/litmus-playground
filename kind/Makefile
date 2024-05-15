REGISTRY := localhost:5000
ORDER_API_DATABASE_IMAGE_NAME := order-api-database
ORDER_API_DATABASE_IMAGE_TAG := latest
ORDER_API_APP_IMAGE_NAME := order-api-app
ORDER_API_APP_IMAGE_TAG := latest
K6_LOADTEST_IMAGE_NAME := k6
K6_LOADTEST_IMAGE_TAG := latest

# order-api

order-api: order-api-database order-api-migration order-api-app

# order-api-database

order-api-database: order-api-database-k8s

order-api-database-k8s:
	kubectl apply -f ../k8s-manifests/order-api-database/namespace.yaml
	kubectl apply -f ../k8s-manifests/order-api-database/pvc.yaml
	kubectl apply -f ../k8s-manifests/order-api-database/deployment.yaml
	kubectl apply -f ../k8s-manifests/order-api-database/service.yaml

# order-api-migration

order-api-migration: order-api-migration-docker order-api-migration-k8s

order-api-migration-docker:
	docker build --no-cache -t $(ORDER_API_DATABASE_IMAGE_NAME) ../k8s-deployments/order-api-migration
	docker tag $(ORDER_API_DATABASE_IMAGE_NAME) $(REGISTRY)/$(ORDER_API_DATABASE_IMAGE_NAME):$(ORDER_API_DATABASE_IMAGE_TAG)
	kind load docker-image $(REGISTRY)/$(ORDER_API_DATABASE_IMAGE_NAME):$(ORDER_API_DATABASE_IMAGE_TAG)

order-api-migration-k8s:
	kubectl apply --force -f ../k8s-manifests/order-api-migration/namespace.yaml
	REGISTRY=$(REGISTRY) IMAGE=$(ORDER_API_DATABASE_IMAGE_NAME) TAG=$(ORDER_API_DATABASE_IMAGE_TAG) envsubst < ../k8s-manifests/order-api-migration/pod.yaml | kubectl delete --force -f - || true
	REGISTRY=$(REGISTRY) IMAGE=$(ORDER_API_DATABASE_IMAGE_NAME) TAG=$(ORDER_API_DATABASE_IMAGE_TAG) envsubst < ../k8s-manifests/order-api-migration/pod.yaml | kubectl apply --force -f -

# order-api-app

order-api-app: order-api-app-docker order-api-app-k8s

order-api-app-docker:
	docker build --no-cache -t $(ORDER_API_APP_IMAGE_NAME) ../k8s-deployments/order-api-app
	docker tag $(ORDER_API_APP_IMAGE_NAME) $(REGISTRY)/$(ORDER_API_APP_IMAGE_NAME):$(ORDER_API_APP_IMAGE_TAG)
	kind load docker-image $(REGISTRY)/$(ORDER_API_APP_IMAGE_NAME):$(ORDER_API_APP_IMAGE_TAG)

order-api-app-k8s:
	kubectl apply --force -f ../k8s-manifests/order-api-app/namespace.yaml
	REGISTRY=$(REGISTRY) IMAGE=$(ORDER_API_APP_IMAGE_NAME) TAG=$(ORDER_API_APP_IMAGE_TAG) envsubst < ../k8s-manifests/order-api-app/deployment.yaml | kubectl delete --force -f - || true
	REGISTRY=$(REGISTRY) IMAGE=$(ORDER_API_APP_IMAGE_NAME) TAG=$(ORDER_API_APP_IMAGE_TAG) envsubst < ../k8s-manifests/order-api-app/deployment.yaml | kubectl apply --force -f -
	kubectl apply --force -f ../k8s-manifests/order-api-app/service.yaml

# k6-loadtest

k6-loadtest: k6-loadtest-docker k6-loadtest-k8s

k6-loadtest-docker:
	docker build --no-cache -t $(K6_LOADTEST_IMAGE_NAME) ../k8s-deployments/k6-loadtest
	docker tag $(K6_LOADTEST_IMAGE_NAME) $(REGISTRY)/$(K6_LOADTEST_IMAGE_NAME):$(K6_LOADTEST_IMAGE_TAG)
	kind load docker-image $(REGISTRY)/$(K6_LOADTEST_IMAGE_NAME):$(K6_LOADTEST_IMAGE_TAG)

k6-loadtest-k8s:
	kubectl apply --force -f ../k8s-manifests/k6-loadtest/namespace.yaml
	REGISTRY=$(REGISTRY) IMAGE=$(K6_LOADTEST_IMAGE_NAME) TAG=$(K6_LOADTEST_IMAGE_TAG) envsubst < ../k8s-manifests/k6-loadtest/deployment.yaml | kubectl delete --force -f - || true
	REGISTRY=$(REGISTRY) IMAGE=$(K6_LOADTEST_IMAGE_NAME) TAG=$(K6_LOADTEST_IMAGE_TAG) envsubst < ../k8s-manifests/k6-loadtest/deployment.yaml | kubectl apply --force -f -
	kubectl apply --force -f ../k8s-manifests/k6-loadtest/service.yaml

# k6-probe

k6-probe: k6-probe-docker

k6-probe-docker:
	docker build --no-cache -t k6-probe ../k6-probe
	docker tag k6-probe localhost:5000/k6-probe:latest
	kind load docker-image localhost:5000/k6-probe:latest
