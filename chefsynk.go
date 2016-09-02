package main

import (
	"encoding/json"
	"fmt"
	"github.com/marpaia/chef-golang"
	"os"
)

// Query Example: "role:spark_sparklecrunch_worker AND chef_environment:vungle_legacy"

var query = os.Getenv("SEARCH_QUERY")

func main() {
	c, err := chef.Connect("knife.rb")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	c.SSLNoVerify = true
	var n chef.Node
	s, err := c.Search("node", query)
	for _, node := range s.Rows {
		err := json.Unmarshal(node, &n)
		if err != nil {
			// Handle
		}
		fmt.Println(n.Info.IPAddress)
	}

}
