package config

import (
	"flag"
	"fmt"
)

type options struct {
	ProxyPort     int
	DashboardPort int
	FrontendPath  string
	MongoUrl      string
	MongoDbName   string
}

func (o *options) GetProxyPortString() string {
	return fmt.Sprintf(":%d", o.ProxyPort)
}

func (o *options) GetDashboardPortString() string {
	return fmt.Sprintf(":%d", o.DashboardPort)
}

var (
	Options *options
)

func init() {
	Options = &options{}
	flag.IntVar(&Options.ProxyPort, "proxyPort", 4000, "Port on which to run the api proxy server")
	flag.IntVar(&Options.DashboardPort, "dashboardPort", 8080, "Port on which to run the dashboard webapp")
	flag.StringVar(&Options.FrontendPath, "frontendPath", "public", "Folder containing the webapp static assets")
	flag.StringVar(&Options.MongoUrl, "mongoUrl", "localhost:27017", "URL to a mongo server")
	flag.StringVar(&Options.MongoDbName, "mongoDbName", "apibutler", "Name of the mongo db to use")
	flag.Parse()
}
