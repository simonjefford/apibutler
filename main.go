package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"fourth.com/apibutler/apiproxyserver"
	"fourth.com/apibutler/config"
	"fourth.com/apibutler/dashboard"
	"fourth.com/apibutler/metadata"
)

func startProxyServer(server apiproxyserver.APIProxyServer) {
	port := config.Options.GetProxyPortString()
	log.Println("Running proxy on", port)
	log.Fatalln(http.ListenAndServe(port, server))
}

func startDashboardServer(proxy apiproxyserver.APIProxyServer, storage metadata.ApiStorage) {
	path := config.Options.FrontendPath
	port := config.Options.GetDashboardPortString()
	server := dashboard.NewDashboardServer(path, proxy, storage)
	log.Println("Running dashboard on", port)
	log.Fatalln(http.ListenAndServe(port, server))
}

func main() {
	// use all available cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	apps := metadata.GetApplicationsTable()

	apiStore := metadata.NewMongoApiStoreFromConfig()

	apis, err := apiStore.Apis()
	if err != nil {
		log.Fatalln(err)
	}

	server := apiproxyserver.NewAPIProxyServer(apps, apis)

	go startProxyServer(server)
	go startDashboardServer(server, apiStore)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
