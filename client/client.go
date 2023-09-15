package main

import (
	"context"
	"flag"
	"fmt"
	"galt-etcd-client/grains"
	"galt-etcd-client/utils"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var DEBUG bool = false
var MINION_BASE string = "/minion"

func init() {
	utils.MinionGrains = grains.NewGrains()
	flag.BoolVar(&DEBUG, "debug", false, "Enable debugging")
	flag.StringVar(&utils.SearchString, "g", "", "Return grain")
	flag.Parse()

	Config = utils.NewConfigFile()
	Config.ConfigInitialize()
	MINION_BASE = fmt.Sprintf("%s/%s", MINION_BASE, Config.Minion.Name)
}

var Config *utils.ConfigFile

func debug(text ...string) {
	if DEBUG {
		fmt.Printf("DEBUG %s\n", strings.Join(text, " "))
	}
}

func getMasterHosts() []string {
	var tmp []string
	for x := range Config.Master.Host {
		tmp = append(tmp, Config.Master.Host[x])
	}
	return tmp
}

func main() {
	// expect dial time-out on ipv4 blackhole
	// _, err := clientv3.New(clientv3.Config{
	// 	Endpoints:   []string{"http://254.0.0.1:12345"},
	// 	DialTimeout: 2 * time.Second,
	// })

	// etcd clientv3 >= v3.2.10, grpc/grpc-go >= v1.7.3
	// if err == context.DeadlineExceeded {
	// 	panic(err)
	// }

	cli, err := clientv3.New(clientv3.Config{
		// Endpoints:   []string{"localhost:2379", "localhost:2379", "localhost:2379"},
		Endpoints:           getMasterHosts(),
		DialTimeout:         5 * time.Second,
		PermitWithoutStream: true,
		DialKeepAliveTime:   time.Duration(time.Duration(Config.Master.KeeplivePing * 1000)),
	})
	if err != nil {
		panic(err)
	}

	defer cli.Close()
	fmt.Println("Connecting to server")

	cluster, _ := cli.Cluster.MemberList(cli.Ctx())
	fmt.Println("etcd: ", cluster)
	utils.MinionGrains.Update()
	fmt.Println(utils.MinionGrains.ToJSON())

	APPLICATION := make(chan string)
	go watchForCmdRun(APPLICATION, cli, MINION_BASE)
	// presp, err := cli.KV.Delete(context.Background(), fmt.Sprintf("foo%d", x))
	// testEtcKV(cli)
	go keepalive(APPLICATION, cli, fmt.Sprintf("%s/connected", MINION_BASE))
	time.Sleep(time.Second * 30)
	APPLICATION <- "quit"
}

func keepalive(APPLICATION chan string, cli *clientv3.Client, key string) {
	for x := 0; x < 100; x += 0 {
		select {
		case applicationAction := <-APPLICATION:
			if strings.Contains(applicationAction, "quit") {
				fmt.Println("quitting go routine watchKey")
				return
			}
		default:
			presp, err := cli.KV.Put(context.Background(), key, "true")

			if err != nil {
				panic(err)
			}
			_ = presp
		}
		time.Sleep(time.Second * 1)
	}
}
func watchForCmdRun(APPLICATION chan string, cli *clientv3.Client, key string) {
	fmt.Println("Watching for server requests on:", key)
	sub := cli.Watch(context.Background(), fmt.Sprintf("%s/request", key), clientv3.WithProgressNotify())
	// wresp := <-sub
	for x := 0; x < 100; x += 0 {
		select {
		case wresp := <-sub:
			for _, ev := range wresp.Events {
				fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				if strings.Contains("cmd.run", string(ev.Kv.Value)) {
					fmt.Printf("Reading the %s/request_content\n", fmt.Sprintf("%s/request_content", key))
					resp, err := cli.KV.Get(context.Background(), fmt.Sprintf("%s/request_content", key), clientv3.WithPrefix())
					if err != nil {
						panic(err)
					}
					requestContent := ""
					for _, ev := range resp.Kvs {
						fmt.Printf("%s : %s\n", ev.Key, ev.Value)
						requestContent = string(ev.Value)
					}
					fmt.Printf("Got the following content to work with: %s\n", requestContent)
				}
			}
		case applicationAction := <-APPLICATION:
			if strings.Contains(applicationAction, "quit") {
				fmt.Println("quitting go routine watchKey")
				return
			}
		default:
			time.Sleep(time.Second * 1)
		}
	}

	fmt.Println("- end watching key:", key)
}

func testEtcKV(cli *clientv3.Client) {
	fmt.Println("Adding 1 kv's")
	for x := 0; x < 1000; x++ {
		presp, err := cli.KV.Put(context.Background(), fmt.Sprintf("/foo_%d", x), fmt.Sprintf("value_%d", x))

		if err != nil {
			panic(err)
		}
		_ = presp
	}

	fmt.Println("Reading the /foo* entries")
	resp, err := cli.KV.Get(context.Background(), "/foo", clientv3.WithPrefix(), clientv3.WithLimit(3))
	if err != nil {
		panic(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	fmt.Println("done")
}
