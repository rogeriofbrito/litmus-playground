# order-api

## Run database

```bash
docker-compose up -d
```

## Run liquibase update

```bash
docker run --rm --network host -v ./migration/changelog:/liquibase/changelog liquibase/liquibase:4.27 --url=jdbc:postgresql://localhost:5432/orderapi?currentSchema=public --changeLogFile=./changelog/changelog.yaml --username=orderapi --password=orderapi update
```

## VS Code launch configuration

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/src",
            "env": {
                "DATABASE_USER": "orderapi",
                "DATABASE_PASSWORD": "orderapi",
                "DATABASE_HOST": "localhost",
                "DATABASE_PORT": "5432",
                "DATABASE_NAME": "orderapi",
                "PORT": "8080"
            }
        }
    ]
}
```

## Run k6

```bash
docker run --rm --network host -e ORDER_HOST=localhost -e ORDER_PORT=8080 -e K6_WEB_DASHBOARD=true -p 5665:5665 -v "$PWD/loadtest:/loadtest" grafana/k6 run /loadtest/script.js --vus 20 --duration 600s
```
