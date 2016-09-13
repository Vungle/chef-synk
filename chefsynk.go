package main

import (
	"encoding/json"
	"fmt"
	"github.com/bradfitz/slice"
	"github.com/marpaia/chef-golang"
	"github.com/zorkian/go-datadog-api"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
)

// Query Example: "role:spark_sparklecrunch_worker AND chef_environment:vungle_legacy"

var query = os.Getenv("SEARCH_QUERY")
var endClient client.EndpointsInterface
var ddog *datadog.Client
var ddog_app_key = os.Getenv("DDOG_APP_KEY")
var ddog_api_key = os.Getenv("DDOG_API_KEY")

func main() {
	// Set sleep interval
	interval, _ := strconv.ParseInt(os.Getenv("INTERVAL"), 10, 64)
	// Connect to chef
	c, err := chef.Connect(os.Getenv("KNIFE_PATH"))
	if err != nil {
		log.Printf("Error:", err)
		os.Exit(1)
	}
	c.SSLNoVerify = true
	// Connect to Kube API
	if kubeClient, err := client.NewInCluster(); err != nil {
		log.Printf("Failed to create client: %v.\n", err)
	} else {
		endClient = kubeClient.Endpoints(os.Getenv("KUBE_NAMESPACE"))
	}
	// Connect to Datadog
	if ddog_api_key != "" && ddog_app_key != "" {
		ddog = datadog.NewClient(ddog_api_key, ddog_app_key)
	}
	// Start Main Loop
	for {
		var chef_nodes = []api.EndpointAddress{}
		var n chef.Node
		// Connect to Chef Server
		s, err := c.Search("node", query)
		if err != nil {
			log.Printf("Failed to connect to chef")
			// Send Event to Data Dog if Enabled
			if ddog != nil {
				event := datadog.Event{
					Title:     "Chef Synk Error:",
					Text:      "Failed to connect to chef",
					AlertType: "error",
					Host:      os.Getenv("HOSTNAME"),
				}
				de, _ := ddog.PostEvent(&event)
				log.Printf("%v", de)
			}
			os.Exit(1)
		}
		// Store Chef Results Like Kube API Formant
		for _, node := range s.Rows {
			err := json.Unmarshal(node, &n)
			if err != nil {
				// Handle
				log.Printf("Failed to unmarshall json")
				os.Exit(1)
			}
			// Get Chef Nodes from Search
			chef_node := api.EndpointAddress{IP: n.Info.IPAddress}
			chef_nodes = append(chef_nodes, chef_node)
		}
		e, err := endClient.Get(os.Getenv("KUBE_ENDPOINT"))
		if err != nil {
			// Handle
			log.Printf("Failed to get kube endpoint")
			os.Exit(1)
		}
		kube_ep_addresses := e.Subsets[0].Addresses
		// Sort before comparison
		slice.Sort(chef_nodes[:], func(i, j int) bool {
			return chef_nodes[i].IP < chef_nodes[j].IP
		})
		slice.Sort(kube_ep_addresses[:], func(i, j int) bool {
			return kube_ep_addresses[i].IP < kube_ep_addresses[j].IP
		})
		// Check if Chef Matches Kube
		if reflect.DeepEqual(chef_nodes, kube_ep_addresses) {
			log.Printf("Kube is synked with Chef!")
		} else if len(chef_nodes) < 1 {
			log.Printf("Chef search returned no nodes: %v", query)
		} else if len(kube_ep_addresses) < 1 {
			log.Printf("Kube endpoint has no addresses: %v", e)
		} else if reflect.DeepEqual(api.EndpointAddress{}, chef_nodes) {
			log.Printf("Can't Update With No Data %v", kube_ep_addresses)
		} else { // All Error Conditions Avoided, Proceed With Update:
			log.Printf("Kube is not in sync with Chef! \nChef Says:\n %v \n\n Kube Says:\n %v \n\n", chef_nodes, kube_ep_addresses)
			e.Subsets[0].Addresses = chef_nodes
			new_eps, err := endClient.Update(e)
			if err != nil {
				log.Printf("Failed to update endpoint:\n %v", err)
				if ddog != nil {
					text := fmt.Sprintf("Failed to Update Endpoint: %v\nFor Endpoint: %v\n", query, os.Getenv("KUBE_ENDPOINT"))
					event := datadog.Event{
						Title:     "Chef Synk Fail:",
						Text:      text,
						AlertType: "error",
						Host:      os.Getenv("HOSTNAME"),
					}
					de, _ := ddog.PostEvent(&event)
					log.Printf("%v", de)
				}
			} else {
				log.Printf("Update suceeded:\n %v", new_eps)
				if ddog != nil {
					text := fmt.Sprintf("Updated Nodes From Search: %v\nFor Endpoint: %v\nNum of Endpoints: %v", query, os.Getenv("KUBE_ENDPOINT"), len(e.Subsets[0].Addresses))
					event := datadog.Event{
						Title:     "Chef Synking Success:",
						Text:      text,
						AlertType: "info",
						Host:      os.Getenv("HOSTNAME"),
					}
					de, _ := ddog.PostEvent(&event)
					log.Printf("%v", de)
				}
			}
		}
		// Wait for interval before starting ( this is for backing off on failure )
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
