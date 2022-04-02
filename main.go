package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"yrfs-exporter/collector"
)

func main() {
	listenAddress := flag.String("listen-address", ":2211", "The listen address of yrfs-exporter.")
	agentServer := flag.String("agent-server", "0.0.0.0:8000", "The yrfs-agent address.")
	namespace := flag.String("namespace", "yrfs", "The name prefix of the metrics.")
	flag.Parse()
	// 将metrics注册到DefaultRegister
	clusterMetrics := collector.NewClusterMetrics(*namespace, *agentServer)
	prometheus.MustRegister(clusterMetrics)

	// 去掉go程序的采集信息
	prometheus.Unregister(prometheus.NewGoCollector())
	prometheus.Unregister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(`
<html>
<body>
<h1>YRCloudfile Exporter</h1>
<p><a href='/metrics'>Metrics</a></p>
</body>
</html>
`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
