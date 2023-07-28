# Keycloak

```bash
kubectl create namespace tools --dry-run=client -o yaml | kubectl apply -f -

kubectl create secret generic keycloak --from-literal=PASSWORD='123456' -n tools

kubectl create secret generic postgres --from-literal=postgres-password='123456' --from-literal=ADMIN_PASSWORD='123456' --from-literal=USER_PASSWORD='123456' --from-literal=password='123456' -n tools

helm install keycloak --values auth/keycloak-helm-values.yaml oci://registry-1.docker.io/bitnamicharts/keycloak -n tools

kubectl apply -f k8s/keycloak-ingress.yaml -n tools
```