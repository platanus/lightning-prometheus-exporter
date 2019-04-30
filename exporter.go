package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/lightningnetwork/lnd/lncfg"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/macaroons"
	"github.com/platanus/lightning-prometheus-exporter/client"
	"github.com/platanus/lightning-prometheus-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	macaroon "gopkg.in/macaroon.v2"
)

func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("Environment Variable value for %s must be a boolean", key)
	}
	return b
}

var (
	// Set during go build
	version   string
	gitCommit string

	// maxMsgRecvSize is the largest message our client will receive. We
	// set this to ~50Mb atm.
	maxMsgRecvSize = grpc.MaxCallRecvMsgSize(1 * 1024 * 1024 * 50)

	// Defaults values
	defaultNamespace     = getEnv("NAMESPACE", "lnd")
	defaultListenAddress = getEnv("LISTEN_ADDRESS", ":9113")
	defaultMetricsPath   = getEnv("TELEMETRY_PATH", "/metrics")
	defaultRPCHost       = getEnv("RPC_HOST", "localhost")
	defaultRPCPort       = getEnv("RPC_PORT", "10009")
	defaultTLSCertPath   = getEnv("TLS_CERT_PATH", "/root/.lnd")
	defaultMacaroonPath  = getEnv("MACAROON_PATH", "")
	defaultGoMetrics, _  = strconv.ParseBool(getEnv("GO_METRICS", "false"))

	// Command-line flags
	namespace = flag.String("namespace", defaultNamespace,
		"The namespace or prefix to use in the exported metrics. The default value can be overwritten by NAMESPACE environment variable.")
	listenAddr = flag.String("web.listen-address", defaultListenAddress,
		"An address to listen on for web interface and telemetry. The default value can be overwritten by LISTEN_ADDRESS environment variable.")
	metricsPath = flag.String("web.telemetry-path", defaultMetricsPath,
		"A path under which to expose metrics. The default value can be overwritten by TELEMETRY_PATH environment variable.")
	rpcHost = flag.String("rpc.host", defaultRPCHost,
		"Lightning node RPC host. The default value can be overwritten by RPC_HOST environment variable.")
	rpcPort = flag.String("rpc.port", defaultRPCPort,
		"Lightning node RPC port. The default value can be overwritten by RPC_PORT environment variable.")
	tlsCertPath = flag.String("lnd.tls-cert-path", defaultTLSCertPath,
		"The path to the tls certificate. The default value can be overwritten by TLS_CERT_PATH environment variable.")
	macaroonPath = flag.String("lnd.macaroon-path", defaultMacaroonPath,
		"The path to the read only macaroon. The default value can be overwritten by MACAROON_PATH environment variable.")
	goMetrics = flag.Bool("go-metrics", defaultGoMetrics,
		"Enable process and go metrics from go client library. The default value can be overwritten by GO_METRICS environmental variable.")
)

func main() {
	flag.Parse()

	log.Printf("Starting Lightning Prometheus Exporter Version=%v GitCommit=%v", version, gitCommit)

	var connCfg = getClientConn()
	rpcclient := lnrpc.NewLightningClient(connCfg)

	client, err := client.NewLightningClient(rpcclient)
	if err != nil {
		log.Fatalf("Could not create Lightning Rpc Client: %v", err)
	}

	// registry
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector.NewLightningCollector(client, *namespace))

	if *goMetrics {
		registry.MustRegister(prometheus.NewGoCollector())
		registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	}

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Lightning Exporter</title></head>
			<body>
			<h1>Lightning Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func getClientConn() *grpc.ClientConn {
	// Load the specified TLS certificate and build transport credentials
	// with it.
	creds, err := credentials.NewClientTLSFromFile(*tlsCertPath, "")
	if err != nil {
		log.Fatalf("could not find TLS certificate: %v", err)
	}

	// Create a dial options array.
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	// Load the specified macaroon file.
	macBytes, err := ioutil.ReadFile(*macaroonPath)
	if err != nil {
		log.Fatalf("could not find Macaroon: %v", err)
	}

	mac := &macaroon.Macaroon{}
	if err = mac.UnmarshalBinary(macBytes); err != nil {
		log.Fatalf("unable to decode macaroon: %v", err)
	}

	// Now we append the macaroon credentials to the dial options.
	cred := macaroons.NewMacaroonCredential(mac)
	opts = append(opts, grpc.WithPerRPCCredentials(cred))

	// We need to use a custom dialer so we can also connect to unix sockets
	// and not just TCP addresses.
	genericDialer := lncfg.ClientAddressDialer(*rpcPort)
	opts = append(opts, grpc.WithDialer(genericDialer))
	opts = append(opts, grpc.WithDefaultCallOptions(maxMsgRecvSize))

	conn, err := grpc.Dial(*rpcHost, opts...)
	if err != nil {
		log.Fatalf("unable to connect to RPC server: %v", err)
	}

	return conn
}
