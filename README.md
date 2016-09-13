# chef-synk
[![Circle CI](https://circleci.com/gh/Vungle/chef-synk)](https://circleci.com/gh/Vungle/docker-chef-synk/tree/master) [![Docker Hub](https://img.shields.io/badge/docker-ready-blue.svg)](https://registry.hub.docker.com/u/vungle/chef-synk/) [![Docker Pulls](https://img.shields.io/docker/pulls/vungle/chef-synk.svg)](https://registry.hub.docker.com/u/vungle/chef-synk/)

To Build
``` make build ```

To run

```
sudo su - 
mkdir -p /mnt/opt
docker run -d -p 443:443 -v /mnt/opt:/var/opt -v /etc/chef-server:/etc/chef-server vungle/chef-server
```


Add Chef Nodes as Endpoints for Kubernetes External Services

## Development

1. Copy `knife.rb` and `vungle.pem` to your local directory (ignored by git) from `/vungle/chef-repo/.chef/`
2. Make sure you are in the minikube context: `kubectl config use-context minikube` and that its running `minikube status`
3. `kubectl create secret generic chef-secrets --from-file vungle.pem --from-file vungle.pem`
4. Validate you have an endpoint available in your minikube context: `k get ep` (that isn't kubernetes)
4. Modify `KUBE_NAMESPACE`, `KUBE_ENDPOINT`, and `KNIFE_SEARCH_QUERY` as needed.
4. Deploy: `kubectl create -f deployment.yaml`
