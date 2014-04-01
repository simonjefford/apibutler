package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"fourth.com/apibutler/apiproxyserver"
	"fourth.com/apibutler/dashboard"
	"fourth.com/apibutler/metadata"
)

type options struct {
	proxyPort     int
	dashboardPort int
	frontendPath  string
}

var (
	opts options
)

func (o options) proxyPortString() string {
	return fmt.Sprintf(":%d", o.proxyPort)
}

func (o options) dashboardPortString() string {
	return fmt.Sprintf(":%d", o.dashboardPort)
}

func init() {
	flag.IntVar(&opts.proxyPort, "proxyPort", 4000, "Port on which to run the rate limiting proxy")
	flag.IntVar(&opts.dashboardPort, "dashboardPort", 8080, "Port on which to run the dashboard webapp")
	flag.StringVar(&opts.frontendPath, "frontendPath", "public", "Folder containing the webapp static assets")
	flag.Parse()
}

func startProxyServer() {
	apps := metadata.GetApplicationsTable()
	apiStore, _ := metadata.GetApiStore()
	apis, _ := apiStore.Apis()
	server := apiproxyserver.NewAPIProxyServer(apps, apis)

	log.Println("Running proxy on", opts.proxyPortString())
	log.Fatalln(http.ListenAndServe(opts.proxyPortString(), server))
}

func startDashboardServer() {
	server := dashboard.NewDashboardServer(opts.frontendPath)
	log.Println("Running dashboard on", opts.dashboardPortString())
	log.Fatalln(http.ListenAndServe(opts.dashboardPortString(), server))
}

func main() {
	// use all available cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	go startProxyServer()
	go startDashboardServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
