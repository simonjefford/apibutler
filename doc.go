// Command apibutler is the main binary for the Fourth API routing layer
//
// Introduction
//
// apibutler is a transparent HTTP proxy that provides a management and routing layer for
// Fourth's REST API endpoints. It also acts as a host for a "thick-client" web application
// for administering the proxy.
//
// Command line flags
//
// -proxyPort - the port on which to run the API proxy server
//
// -dashboardPort - the port on whith to run the dashboard admin app
//
// -frontendPath - the path from which to serve the dashboard admin app
package main
