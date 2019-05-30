# kubernetes-recipes
A collection of Kubernetes recipes.

## Services

- PiHole DNS Server
- Docker Registry
- CertBot
- Ingress
- Plex Server (ingressed/nodeport)
- Volumes NFS (broken)

## Usage

``` sh
kubectl get all -n hebron           # View all elements for namespace 'hebron'
kubectl create -f k8s/[SERVICE].yaml      # Create all elements defined in yaml spec
kubectl delete -f k8s/[SERVICE].yaml      # Delete all elements defines in yaml spec
kubectl scale deployment [SERVICE] --replicas=0 -n NAMESPACE # shutdown service by removing all pods (0)
kubectl scale deployment [SERVICE] --replicas=1 -n NAMESPACE # start service by scaling pods (1)
```