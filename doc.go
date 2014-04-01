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
//
// Startup
//
// The startup process goes like this:
//
// 1. Parse any command line flags.
//
// 2. Startup the dashboard server on the configured port, serving files from the configured
// directory
//
// 3. Startup the proxy server on the configured port.
//
// 4. Wait for SIGINT
package main
