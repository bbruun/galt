package utils

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestMinions() {
	/*
		This function will find all /minion/* entries and will then
		send random "crap" to them to execute to test functionality.
	*/
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:           []string{"127.0.0.1:2379"},
		DialTimeout:         1 * time.Second,
		PermitWithoutStream: true,
	})
	if err != nil {
		panic(err)
	}

	defer cli.Close()
	cluster, _ := cli.Cluster.MemberList(cli.Ctx())
	fmt.Printf("cluster: %+v\n", cluster)
	fmt.Println("Reading the /minion* entries")
	for x := 0; x < 100; x += 0 {
		resp, err := cli.KV.Get(context.Background(), "/minion/", clientv3.WithPrefix())
		if err != nil {
			panic(err)
		}
		for _, ev := range resp.Kvs {
			key := string(ev.Key)
			val := string(ev.Value)
			if strings.Contains(key, "/connected") {
				log.Printf("(/) - %s : %s\n", ev.Key, val)
			} else {
				log.Printf("(x) - %s : %s\n", ev.Key, val)
			}
		}
		time.Sleep(time.Second * 1)
	}
	fmt.Println("done")
}
