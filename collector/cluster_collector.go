package collector

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"log"
	pb "yrfs-exporter/proto"
	"yrfs-exporter/yrfs"
)

// 定义ClusterMetrics
type ClusterMetrics struct {
	metrics         map[string]*prometheus.Desc
	yrfsAgentServer string
}

// 实例化Metrics
func newGlobalMetrics(namespace string, metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, labels, nil)
}

func NewClusterMetrics(namespace string, yrfsAgentServer string) *ClusterMetrics {
	return &ClusterMetrics{
		yrfsAgentServer: yrfsAgentServer,
		metrics: map[string]*prometheus.Desc{
			"total_metas": newGlobalMetrics(namespace, "total_metas", "Number of the MDS", []string{"type"}),
			"online_metas": newGlobalMetrics(namespace, "online_metas", "Number of the online MDS", []string{"type"}),
			"meta_capacity_total": newGlobalMetrics(namespace, "meta_capacity_total", "The capacity of the MDS", []string{"type"}),
			"meta_capacity_used":  newGlobalMetrics(namespace, "meta_capacity_used", "The used capacity of the MDS", []string{"type"}),
		},
	}
}

// Describe 方法用于返回独立的Desc实例
func (c *ClusterMetrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

// 当promhttp.Handler()被调用时，Collect函数被执行
func (c *ClusterMetrics) Collect(ch chan<- prometheus.Metric) {
	// call the grpc server of yrfs-agent
	resp, err := mdsOverview(c.yrfsAgentServer, context.Background())
	if err != nil {
		return
	}

	metaData := make(map[string]uint64)

	// 获取集群的副本数
	metaCopies, _, err := yrfs.GetClustercopies()
	if err != nil {
		log.Fatal("get cluster copies fail: %v", err)
		return
	}

	// 获取元数据服务总数， 在线元数据服务数
	var nodes, onlineNodes uint8
	for _, node := range resp.NodeInfo {
		nodes += 1
		if node.Online == true {
			onlineNodes += 1
		}
	}

	metaData["total_metas"] = uint64(nodes)
	metaData["online_metas"] = uint64(onlineNodes)
	metaData["meta_capacity_total"] = resp.DiskSpaceTotal / uint64(metaCopies)
	metaData["meta_capacity_used"] = resp.DiskSpaceUsed / uint64(metaCopies)
	metaData["meta_capacity_available"] = resp.DiskSpaceFree / uint64(metaCopies)
	metaData["meta_inode_capacity_used"] = resp.InodeSpaceUsed / uint64(metaCopies)

	for k, currentValue := range metaData {
		ch <- prometheus.MustNewConstMetric(c.metrics[k], prometheus.GaugeValue, float64(currentValue), "cluster")
	}
}

func mdsOverview(agentServer string, ctx context.Context) (ret *pb.MdsOverviewRet, err error) {
	conn, err := grpc.Dial(agentServer, grpc.WithInsecure())
	if err != nil {
		fmt.Println("grpc server connect failed")
		return
	}
	defer conn.Close()
	c := pb.NewAgentClient(conn)
	resp, err := c.MdsOverview(ctx, &pb.MdsOverviewPara{})
	if err != nil {
		fmt.Println("could not greet:")
	}

	return resp, nil
}
