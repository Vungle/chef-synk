package main

import (
	"encoding/json"
	"fmt"
	"github.com/marpaia/chef-golang"
	_ "k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"os"
	"strconv"
	"time"
)

// Query Example: "role:spark_sparklecrunch_worker AND chef_environment:vungle_legacy"

var query = os.Getenv("SEARCH_QUERY")

func main() {
	interval, _ := strconv.ParseInt(os.Getenv("INTERVAL"), 10, 64)
	c, err := chef.Connect(os.Getenv("KNIFE_PATH"))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	c.SSLNoVerify = true
	for {
		var n chef.Node
		s, _ := c.Search("node", query)
		for _, node := range s.Rows {
			err := json.Unmarshal(node, &n)
			if err != nil {
				// Handle
				fmt.Println("Failed to unmarshall json")
			}
			fmt.Println(n.Info.IPAddress)
		}
		ingress()
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func ingress() {
	var endClient client.EndpointsInterface
	if kubeClient, err := client.NewInCluster(); err != nil {
		fmt.Printf("Failed to create client: %v.\n", err)
	} else {
		endClient = kubeClient.Endpoints(os.Getenv("KUBE_NAMESPACE"))
	}
	e, _ := endClient.Get(os.Getenv("KUBE_ENDPOINT"))
	for _, s := range e.Subsets {
		fmt.Printf("%v", s.Addresses)
	}
}
