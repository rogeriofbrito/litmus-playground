# order-api

## Run database

```bash
docker-compose up -d
```

## Run liquibase update

```bash
docker run --rm --network host -v ./migration/changelog:/liquibase/changelog liquibase/liquibase:4.27 --url=jdbc:postgresql://localhost:5432/orderapi?currentSchema=public --changeLogFile=./changelog/changelog.yaml --username=orderapi --password=orderapi update
```

## Run k6

```bash
k6 run ../k6-loadtest/script.js -e ORDER_HOST=localhost -e ORDER_PORT=8080 --vus 20 --duration 600s -e K6_WEB_DASHBOARD=true
```
