# chef-synk
Add Chef Nodes as Endpoints for Kubernetes External Services

## Development

1. Copy `knife.rb` and `vungle.pem` to your local directory (ignored by git) from `/vungle/chef-repo/.chef/`
2. Make sure you are in the minikube context: `kubectl config use-context minikube` and that its running `minikube status`
3. `kubectl create secret generic chef-secrets --from-file vungle.pem --from-file vungle.pem`
4. Validate you have an endpoint available in your minikube context: `k get ep` (that isn't kubernetes)
4. Modify `KUBE_NAMESPACE`, `KUBE_ENDPOINT`, and `KNIFE_SEARCH_QUERY` as needed.
4. Deploy: `kubectl create -f deployment.yaml`
