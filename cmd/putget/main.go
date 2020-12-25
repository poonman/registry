package main

import (
	"fmt"
	"github.com/poonman/registry/registry"
	"github.com/poonman/registry/registry/etcd"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//endpoints := []string{
	//	"192.168.81.51:2379",
	//	"192.168.81.51:3379",
	//	"192.168.81.51:4379",
	//}

	endpoints := []string {
		"slb-pre-server.ifere.com:31334",
	}

	//endpoints := []string {
	//	"infra-etcd-cluster-1.infra-etcd-cluster.demo:",
	//	"infra-etcd-cluster-2.infra-etcd-cluster.demo",
	//	"infra-etcd-cluster-3.infra-etcd-cluster.demo",
	//	"infra-etcd-cluster-4.infra-etcd-cluster.demo",
	//	"infra-etcd-cluster-5.infra-etcd-cluster.demo",
	//	"infra-etcd-cluster-6.infra-etcd-cluster.demo",
	//	"infra-etcd-cluster-7.infra-etcd-cluster.demo",
	//}
	r := etcd.NewRegistry(registry.Addrs(endpoints...))
	err := r.Init()
	if err != nil {
		fmt.Println("err: ", err)
		os.Exit(1)
	}

	for i:=0; i<500; i++ {

		go func(j int) {
			name := fmt.Sprintf("svc-%d", j)
			s := &registry.Service{
				Name:      name,
				Version:   "0.0.0",
				Metadata:  nil,
				Endpoints: nil,
				Nodes:     []*registry.Node{
					{
						Id:       name,
						Address:  "192.168.81.51:9000",
						Metadata: nil,
					},
				},
			}
			err = r.Register(s)
			if err != nil {
				fmt.Println("Failed to register. i: ", j, ". err: ", err)
			}

			fmt.Println("register success. name: ", name)

			w, err := r.Watch(registry.WatchService(name))
			if err != nil {
				fmt.Println("Failed to watch. i: ", j, ". err: ", err)
				return
			}

			res, err := w.Next()
			if err != nil {
				fmt.Println("Failed to watch next. i: ", j, ". err: ", err)
				return
			}

			fmt.Println("action: ", res.Action, ", name: ", res.Service.Name)
		}(i)

	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:

			for i:=0; i<500;i++ {
				name := fmt.Sprintf("svc-%d", i)
				s := &registry.Service{
					Name:      name,
					Version:   "0.0.0",
					Metadata:  nil,
					Endpoints: nil,
					Nodes:     []*registry.Node{
						{
							Id:       name,
							Address:  "192.168.81.51:9000",
							Metadata: nil,
						},
					},
				}
				err = r.Deregister(s)
				if err != nil {
					fmt.Println("Failed to deregister. ", err)
				}
			}

			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
