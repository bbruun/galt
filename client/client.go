package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"galt-etcd-client/grains"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/cameronnewman/go-flatten"
	"github.com/tidwall/gjson"
	"go.etcd.io/etcd/clientv3"
)

var minionGrains *grains.Grains

func init() {
	minionGrains = grains.NewGrains()
}

func getFlattenMap(s string) flatten.Map {
	var jsonBlob map[string]interface{}
	d := json.NewDecoder(bytes.NewReader([]byte(s)))
	d.UseNumber()
	if err := d.Decode(&jsonBlob); err != nil {
		log.Fatal(err)
	}
	// flattenedObject := flatten.Flatten(jsonBlob)
	return flatten.Flatten(jsonBlob)
}

func main() {

	var searchString string
	flag.StringVar(&searchString, "g", "", "Return grain")
	flag.Parse()
	var searchList []string
	if strings.Contains(searchString, " or ") {
		searchList = strings.Split(searchString, " or ")
	}
	// expect dial time-out on ipv4 blackhole
	_, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://254.0.0.1:12345"},
		DialTimeout: 2 * time.Second,
	})

	// etcd clientv3 >= v3.2.10, grpc/grpc-go >= v1.7.3
	if err == context.DeadlineExceeded {
		panic(err)
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:2379", "localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		panic(err)
	}
	defer cli.Close()

	minionGrains.Update()

	// Test output to etcd
	s := minionGrains.ToJSON()
	fmt.Println(s)

	// START SEARCH FOR GRAIN
	// Using "github.com/cameronnewman/go-flatten"
	fmt.Printf("search: %s\n", searchString)
	if searchString != "" {
		fmt.Println("- using \"g\" to find grain")

		flattenedObject := getFlattenMap(s)
		for searchResult := range searchList {
			search := searchList[searchResult]
			if strings.Contains(search, "=") {
				skeyval := strings.Split(search, "=")
				sreg, _ := regexp.Compile(strings.ToLower(skeyval[1]))
				for k, v := range flattenedObject {
					// fmt.Printf(" (%s = %s)\n", k, v)
					if strings.Contains(strings.ToLower(k), skeyval[0]) {
						regsearch := sreg.FindString(strings.ToLower(v))
						if regsearch != "" {
							fmt.Println(k, v)
						}
					}
				}
			} else if strings.Contains(search, ">") || strings.Contains(search, "<") {
				fmt.Println(" - numeric grain search not supported yet")
			} else {
				for k, v := range flattenedObject {
					if strings.Contains(strings.ToLower(k), strings.ToLower(search)) {
						fmt.Println(k, v)
					}
				}
			}
		}
	} else {
		fmt.Println("finding sysinfo.bios.vendor")
		vs := gjson.Get(strings.ToLower(s), "sysinfo.bios.vendor")

		fmt.Println(vs.String())
	}
	// END SEARCH FOR GRAIN
	os.Exit(0)
	fmt.Println("Adding 1 kv's")
	for x := 0; x < 1000; x++ {
		presp, err := cli.KV.Put(context.Background(), fmt.Sprintf("/foo_%d", x), fmt.Sprintf("value_%d", x))
		// presp, err := cli.KV.Delete(context.Background(), fmt.Sprintf("foo%d", x))
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
