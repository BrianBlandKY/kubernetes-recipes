# kubernetes-recipes
A collection of Kubernetes recipes.

## Shawshank Namespace

Internal services for personal usage.

- PiHole DNS Server
- Samba Server
- Plex Server
- Others? (VPN?)

## Development

- Drone Server

## NothingsBland Namespace

- WebServer (Production)
- WebServer (Internal)

## Usage

``` sh
kubectl get all --namespace=shawshank           # View all elements for namespace 'shawshank'
kubectl create -f shawshank/[SERVICE].yaml      # Create all elements defined in yaml spec
kubectl delete -f shawshank/[SERVICE].yaml      # Delete all elements defines in yaml spec
```