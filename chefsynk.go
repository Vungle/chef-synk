package main

import (
	"encoding/json"
	"fmt"
	"github.com/bradfitz/slice"
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
var endClient client.EndpointsInterface

func main() {
	interval, _ := strconv.ParseInt(os.Getenv("INTERVAL"), 10, 64)
	c, err := chef.Connect(os.Getenv("KNIFE_PATH"))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	c.SSLNoVerify = true
	if kubeClient, err := client.NewInCluster(); err != nil {
		fmt.Printf("Failed to create client: %v.\n", err)
	} else {
		endClient = kubeClient.Endpoints(os.Getenv("KUBE_NAMESPACE"))
	}
	// Start Main Loop
	for {
		var chef_nodes = []api.EndpointAddress{}
		var n chef.Node
		// Connect to Chef Server
		s, _ := c.Search("node", query)
		for _, node := range s.Rows {
			err := json.Unmarshal(node, &n)
			if err != nil {
				// Handle
				fmt.Println("Failed to unmarshall json")
				os.Exit(1)
			}
			// Get Chef Nodes from Search
			chef_node := api.EndpointAddress{IP: n.Info.IPAddress}
			chef_nodes = append(chef_nodes, chef_node)
		}
		e, _ := endClient.Get(os.Getenv("KUBE_ENDPOINT"))
		kube_ep_addresses := e.Subsets[0].Addresses
		// Sort before comparison
		slice.Sort(chef_nodes[:], func(i, j int) bool {
			return chef_nodes[i].IP < chef_nodes[j].IP
		})
		slice.Sort(kube_ep_addresses[:], func(i, j int) bool {
			return kube_ep_addresses[i].IP < kube_ep_addresses[j].IP
		})
		// Compare
		if !reflect.DeepEqual(chef_nodes, kube_ep_addresses) {
			fmt.Printf("\nKube is not in sync with Chef! \nChef Says:\n %v \n\n Kube Says:\n %v \n\n", chef_nodes, kube_ep_addresses)
			// Check if Chef Search is Empty
			if !reflect.DeepEqual(api.EndpointAddress{}, chef_nodes) {
				// Modify Kube Response to Match Chef Response
				e.Subsets[0].Addresses = chef_nodes
				new_eps, err := endClient.Update(e)
				if err != nil {
					// Handle
					fmt.Println("Failed to update endpoint:\n %v", err)
					os.Exit(1)
				} else {
					fmt.Printf("Update suceeded:\n %v\n", new_eps)
				}
			} else {
				fmt.Println("Avoiding Updating With No Data")
			}
		} else {
			fmt.Println("Kube is in sync with Chef!")
		}
		// Wait for interval before looping again
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
