# litmus-playgraound

## How to install kubernetes-playground helm chart?

```bash
DATABASE_USER=<database-user>
DATABASE_PASSWORD=<database-password>

helm install kubernetes-playground ./helm \
-n kubernetes-playground \
--create-namespace \
--set order-api.secrets.DATABASE_USER="$(echo -n $DATABASE_USER | base64)" \
--set order-api.secrets.DATABASE_PASSWORD="$(echo -n $DATABASE_PASSWORD | base64)"
```
