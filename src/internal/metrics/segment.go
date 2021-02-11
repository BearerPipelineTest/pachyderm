package metrics

import (
	"fmt"
	"time"

	"github.com/segmentio/analytics-go"
)

const reportingInterval time.Duration = 60 * time.Minute

func newPersistentClient() *analytics.Client {
	c := newSegmentClient()
	c.Interval = reportingInterval
	c.Size = 100
	return c
}

func newSegmentClient() *analytics.Client {
	return analytics.New("hhxbyr7x50w3jtgcwcZUyOFrTf4VNMrD")
}

func reportClusterMetricsToSegment(client *analytics.Client, metrics *Metrics) {
	// We're intentionally ignoring an error here because metrics code is
	// non-critical
	client.Track(&analytics.Track{
		Event:       "cluster.metrics",
		AnonymousId: metrics.ClusterID,
		Properties: map[string]interface{}{
			"PodID":           metrics.PodID,
			"nodes":           metrics.Nodes,
			"version":         metrics.Version,
			"repos":           metrics.Repos,
			"commits":         metrics.Commits,
			"files":           metrics.Files,
			"bytes":           metrics.Bytes,
			"jobs":            metrics.Jobs,
			"pipelines":       metrics.Pipelines,
			"ActivationCode":  metrics.ActivationCode,
			"MaxBranches":     metrics.MaxBranches,
			"PpsSpout":        metrics.PpsSpout,
			"PpsSpoutService": metrics.PpsSpoutService,
			"PpsBuild":        metrics.PpsBuild,
			"CfgEgress":       metrics.CfgEgress,
			"CfgStandby":      metrics.CfgStandby,
			"CfgS3Gateway":    metrics.CfgS3Gateway,
			"CfgServices":     metrics.CfgServices,
			"CfgErrcmd":       metrics.CfgErrcmd,
			"CfgStats":        metrics.CfgStats,
			"CfgTfjob":        metrics.CfgTfjob,
			"InputGroup":      metrics.InputGroup,
			"InputJoin":       metrics.InputJoin,
			"InputCross":      metrics.InputCross,
			"InputUnion":      metrics.InputUnion,
			"InputCron":       metrics.InputCron,
			"InputGit":        metrics.InputGroup,
			"InputOuterJoin":  metrics.InputOuterJoin,
			"InputLazy":       metrics.InputLazy,
			"InputEmptyFiles": metrics.InputEmptyFiles,
			"InputS3":         metrics.InputS3,
			"InputTrigger":    metrics.InputTrigger,
			"ResourceCpu":     metrics.ResourceCpu,
			"ResourceMem":     metrics.ResourceMem,
			"ResourceGpu":     metrics.ResourceGpu,
			"ResourceDisk":    metrics.ResourceDisk,
			"MaxParallelism":  metrics.MaxParallelism,
			"MinParallelism":  metrics.MinParallelism,
		},
	})
}

/*
Segment needs us to identify a user before we report any events for that user.
We have no way of knowing if a user has previously been identified, so we call this
before every `Track()` call containing user data.
*/
func identifyUser(client *analytics.Client, userID string) {
	// We're intentionally ignoring an error here because metrics code is
	// non-critical
	client.Identify(&analytics.Identify{
		UserId: userID,
	})
}

func reportUserMetricsToSegment(client *analytics.Client, userID string, prefix string, action string, value interface{}, clusterID string) {
	identifyUser(client, userID)
	properties := map[string]interface{}{
		"ClusterID": clusterID,
	}
	properties[action] = value
	// We're intentionally ignoring an error here because metrics code is
	// non-critical
	client.Track(&analytics.Track{
		Event:      fmt.Sprintf("%v.usage", prefix),
		UserId:     userID,
		Properties: properties,
	})
}
