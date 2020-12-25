package main

import (
	"flag"
	"fmt"
	"github.com/poonman/registry/registry"
	"github.com/poonman/registry/registry/cache"
	"github.com/poonman/registry/registry/etcd"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	localName = flag.String("local", "svc-a", "")
	remoteName = flag.String("remote", "svc-b", "")
)

func main() {

	flag.Parse()

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

	s := &registry.Service{
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
	}

	err = c.Register(s, registry.RegisterTTL(5*time.Second))

	if err != nil {
		fmt.Println("err: ", err)
		os.Exit(1)
	}


	fmt.Println("local: ", *localName)
	fmt.Println("remote: ", *remoteName)

	go func() {
		for {
			ss, err := c.GetService(*remoteName)
			if err != nil {
				fmt.Println("Failed to get service. ", err)
				time.Sleep(time.Second)
				continue
			}



			for _, s := range ss {
				fmt.Println("Success. ", *s)
			}

			time.Sleep(time.Second)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:

			err = c.Deregister(s)
			if err != nil {
				fmt.Println("Failed to deregister. ", err)
			}
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

