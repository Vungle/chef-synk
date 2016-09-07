package main

import (
	"encoding/json"
	"fmt"
	"github.com/marpaia/chef-golang"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"os"
	"reflect"
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
		var chef_nodes = []api.EndpointAddress{}
		var n chef.Node
		s, _ := c.Search("node", query)
		for _, node := range s.Rows {
			err := json.Unmarshal(node, &n)
			if err != nil {
				// Handle
				fmt.Println("Failed to unmarshall json")
			}
			//fmt.Println(n.Info.IPAddress)
			chef_node := api.EndpointAddress{IP: n.Info.IPAddress}
			chef_nodes = append(chef_nodes, chef_node)
		}
		kube_ep_addresses := endpoints()
		if !reflect.DeepEqual(chef_nodes, kube_ep_addresses) {
			fmt.Printf("Kube is not in sync with Chef! \nChef Says:\n %v \n\n Kube Says:\n %v \n\n", chef_nodes, kube_ep_addresses)
		} else {
			fmt.Println("Kube is in sync with Chef!")
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func endpoints() []api.EndpointAddress {
	var endClient client.EndpointsInterface
	if kubeClient, err := client.NewInCluster(); err != nil {
		fmt.Printf("Failed to create client: %v.\n", err)
	} else {
		endClient = kubeClient.Endpoints(os.Getenv("KUBE_NAMESPACE"))
	}
	e, _ := endClient.Get(os.Getenv("KUBE_ENDPOINT"))
	/*for _, s := range e.Subsets {
		fmt.Printf("%v", s.Addresses)
	}*/
	return e.Subsets[0].Addresses
}
