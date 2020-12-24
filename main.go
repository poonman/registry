package main

import (
	"flag"
	"fmt"
	"github.com/poonman/registry/registry"
	"github.com/poonman/registry/registry/cache"
	"github.com/poonman/registry/registry/etcd"
	"os"
	"time"
)

var (
	localName = flag.String("local", "svc-a", "")
	remoteName = flag.String("remote", "svc-b", "")
)

func main() {
	endpoints := []string{
		"192.168.81.51:2379",
		"192.168.81.51:3379",
		"192.168.81.51:4379",
	}
	r := etcd.NewRegistry(registry.Addrs(endpoints...))
	err := r.Init()
	if err != nil {
		fmt.Println("err: ", err)
		os.Exit(1)
	}

	c := cache.New(r)

	err = c.Register(&registry.Service{
		Name:      *localName,
		Version:   "0.0.0",
		Metadata:  nil,
		Endpoints: nil,
		Nodes:     []*registry.Node{
			{
				Id:       fmt.Sprintf("%s-1", *localName),
				Address:  "192.168.81.51:9000",
				Metadata: nil,
			},
		},
	})

	if err != nil {
		fmt.Println("err: ", err)
		os.Exit(1)
	}

	for {
		ss, err := c.GetService(*remoteName)
		if err != nil {
			fmt.Println("Failed to get service. ", err)
			time.Sleep(time.Second)
			continue
		}

		fmt.Println("get begin...")

		for _, s := range ss {
			fmt.Println(*s)
		}

		fmt.Println("get end.....")

		time.Sleep(time.Second)
	}
}
