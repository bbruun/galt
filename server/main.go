package main

import (
	"flag"
	"fmt"
	"galt-etcd/utils"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"go.etcd.io/etcd/server/v3/embed"
)

const (
	ERROR_NO_BIND_IP = 1
)

func parseEtcdUrls(strs []string) []url.URL {
	urls := make([]url.URL, 0, len(strs))
	for _, str := range strs {
		u, err := url.Parse(str)
		if err != nil {
			log.Printf("Invalid url %s, error: %s", str, err.Error())
			continue
		}
		urls = append(urls, *u)
	}
	return urls
}

// Command line flags
var DEBUG bool = false
var SystemName string = "Galt"
var Storage string = "./etcd-storage"
var BindAddress string = "0.0.0.0"
var EtcdLogLevel string = "error"
var SEND_RANDOM_STUFF bool = false

func debug(msg ...string) {
	if DEBUG {
		log.Printf("DEBUG: %s\n", strings.Join(msg, " "))
	}
}

func init() {

	// Enable DEBUG
	flag.BoolVar(&DEBUG, "debug", false, "Enable debug if any")
	flag.BoolVar(&DEBUG, "d", false, "Enable debug if any")
	// A name for the embedded etcd system
	flag.StringVar(&SystemName, "name", SystemName, "The name of the Galt service - default \"galt\"")
	flag.StringVar(&SystemName, "n", SystemName, "The name of the Galt service - default \"galt\"")
	// Where should the embedded etcd server persist data
	flag.StringVar(&Storage, "storage", Storage, "The storage directory for etcd")
	flag.StringVar(&Storage, "s", Storage, "The storage directory for etcd")
	// What IP should the embedded etcd server bind to
	flag.StringVar(&BindAddress, "bind", BindAddress, "Set the IP to bind the etcd service to - aka the IP clients will connect to on this server")
	flag.StringVar(&BindAddress, "b", BindAddress, "Set the IP to bind the etcd service to - aka the IP clients will connect to on this server")
	// Change embedded etcd loglevel
	flag.StringVar(&EtcdLogLevel, "etcd-loglevel", EtcdLogLevel, "The loglevel for etcd (debug, info, warn, error, panic, or fatal)")

	// Test all minions
	flag.BoolVar(&SEND_RANDOM_STUFF, "test-all-minions", false, "Send random state.sls and cmd.run (non harmful) to all minions")
	flag.Parse()

	debug("Debug is enabled:", strconv.FormatBool(DEBUG))
	debug("Storage set to  :", SystemName)
	debug("Bind address    :", BindAddress)

	if BindAddress == "" {
		fmt.Println("The server needs an IP address to bind to !")
		os.Exit(ERROR_NO_BIND_IP)
	}
	// os.Exit(0)
}

func setupEtcdEmbed() (*embed.Etcd, error) {

	config := embed.NewConfig()

	config.Name = SystemName
	config.Dir = Storage // "/tmp/my-embedded-ectd-cluster"
	config.ListenPeerUrls = parseEtcdUrls([]string{"http://0.0.0.0:2380"})
	config.ListenClientUrls = parseEtcdUrls([]string{"http://0.0.0.0:2379"})
	config.AdvertisePeerUrls = parseEtcdUrls([]string{fmt.Sprintf("http://%s:2380", BindAddress)})
	config.AdvertiseClientUrls = parseEtcdUrls([]string{fmt.Sprintf("http://%s:2380", BindAddress)})
	config.InitialCluster = fmt.Sprintf("%s=http://%s:2380", SystemName, BindAddress)
	config.SelfSignedCertValidity = 2
	config.TlsMinVersion = "TLS1.3"
	config.ClientAutoTLS = true
	config.PeerAutoTLS = true
	config.PeerSelfCert()

	config.LogLevel = EtcdLogLevel

	return embed.StartEtcd(config)
}

func main() {

	etcd, err := setupEtcdEmbed()

	if err != nil {

		log.Fatal(err)

	}
	defer etcd.Close()

	select {
	case <-etcd.Server.ReadyNotify():
		log.Printf("Etcd server is ready!")
	}
	if SEND_RANDOM_STUFF == true {
		fmt.Println("test-all-minions have been passed... please wait...")
		utils.TestMinions()
	}

	err = <-etcd.Err()
	log.Fatal(err)
}
