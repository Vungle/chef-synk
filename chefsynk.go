package main

import (
	"fmt"
	"github.com/marpaia/chef-golang"
	"os"
)

var findNode = os.Getenv("CHEF_SYNK_NODE")

func main() {
	c, err := chef.Connect("knife.rb")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	c.SSLNoVerify = true
	nodes, err := c.GetNodes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// do what you please with the "node" variable which is a map of
	// node names to node URLs
	for node := range nodes {
		fmt.Println(node)
	}

}
