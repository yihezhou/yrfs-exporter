package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"yrfs-exporter/collector"
)

func main() {
	// 将metrics注册到DefaultRegister
	clusterMetrics := collector.NewClusterMetrics("yrfs", "0.0.0.0:8000")
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
	log.Fatal(http.ListenAndServe(":2112", nil))
}
