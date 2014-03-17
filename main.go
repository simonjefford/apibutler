package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"fourth.com/ratelimit/apiproxyserver"
	"fourth.com/ratelimit/dashboard"
	"fourth.com/ratelimit/limiter"
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

func startLimitServer() {
	server := apiproxyserver.NewAPIProxyServer()

	log.Println("Running proxy on", opts.proxyPortString())
	log.Fatalln(http.ListenAndServe(opts.proxyPortString(), server))
}

func startDashboardServer(r limiter.RateLimit) {
	server := dashboard.NewDashboardServer(r, opts.frontendPath)
	log.Println("Running dashboard on", opts.dashboardPortString())
	log.Fatalln(http.ListenAndServe(opts.dashboardPortString(), server))
}

func main() {
	// use all available cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	go startLimitServer()

	// temporary measure
	l, _ := limiter.NewRateLimit()
	go startDashboardServer(l)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
