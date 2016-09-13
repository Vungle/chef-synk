# chef-synk
[![CircleCI](https://circleci.com/gh/Vungle/chef-synk/tree/master.svg?style=svg)](https://circleci.com/gh/Vungle/chef-synk/tree/master) [![Docker Hub](https://img.shields.io/badge/docker-ready-blue.svg)](https://registry.hub.docker.com/u/vungle/chef-synk/) [![Docker Pulls](https://img.shields.io/docker/pulls/vungle/chef-synk.svg)](https://registry.hub.docker.com/u/vungle/chef-synk/)

## Description

Add Chef Nodes as Endpoints for Kubernetes External Services

## Installation

1. Copy `knife.rb` and `vungle.pem` to your local directory (ignored by git) from `/vungle/chef-repo/.chef/`
2. Make sure you are in the minikube context: `kubectl config use-context minikube` and that its running `minikube status`
3. `kubectl create secret generic chef --from-file vungle.pem --from-file knife.rb`
4. Validate you have an endpoint available in your minikube context: `k get ep` (that isn't kubernetes)
4. Modify `KUBE_NAMESPACE`, `KUBE_ENDPOINT`, and `KNIFE_SEARCH_QUERY` as needed.
4. Deploy: `kubectl create -f deployment.yaml`


