package main

import (
	"flag"
	"fmt"
	"fourth.com/ratelimit/apiproxyserver"
	"fourth.com/ratelimit/dashboard"
	"fourth.com/ratelimit/limiter"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
)

type options struct {
	proxyPort     int
	dashboardPort int
	publicPath    string
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
	flag.StringVar(&opts.publicPath, "publicPath", "public", "Folder containing the webapp static assets")
	flag.Parse()
}

func startLimitServer(r *limiter.RateLimit) {
	server := apiproxyserver.NewProxyServer(r)
	log.Println("Running proxy on", opts.proxyPortString())
	log.Fatalln(http.ListenAndServe(opts.proxyPortString(), server))
}

func startDashboardServer(r *limiter.RateLimit) {
	server := dashboard.NewDashboardServer(r, opts.publicPath)
	log.Println("Running dashboard on", opts.dashboardPortString())
	log.Fatalln(http.ListenAndServe(opts.dashboardPortString(), server))
}

func main() {
	// use all available cores
	runtime.GOMAXPROCS(runtime.NumCPU())
	r, err := limiter.NewRateLimit()
	if err != nil {
		log.Fatalln(err)
	}
	go startLimitServer(r)
	go startDashboardServer(r)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
