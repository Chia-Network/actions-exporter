package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	gh "github.com/chia-network/actions-exporter/internal/github"
)

// Collector contains all prometheus metric Descs
type Collector struct {
	orgWorkflowState *prometheus.Desc
	rateLimitLimit   *prometheus.Desc
	rateLimitUsed    *prometheus.Desc
}

// NewCollector constructor function for Collector
func NewCollector() *Collector {
	return &Collector{
		orgWorkflowState: prometheus.NewDesc("github_organization_workflow_state",
			"Shows non-active workflow state for workflows in a GitHub organization.",
			[]string{"owner", "repository", "workflow", "state"}, nil,
		),
		rateLimitLimit: prometheus.NewDesc("github_rate_limit_limit",
			"The maximum number of requests allowed per hour (x-ratelimit-limit).",
			nil, nil,
		),
		rateLimitUsed: prometheus.NewDesc("github_rate_limit_used",
			"The number of requests used in the current rate limit window (x-ratelimit-used).",
			nil, nil,
		),
	}
}

// Describe contains all the prometheus descriptors for this metric collector
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.orgWorkflowState
	ch <- c.rateLimitLimit
	ch <- c.rateLimitUsed
}

// Collect instructs the prometheus client how to collect the metrics for each descriptor
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	t := getOrgWorkflowState()
	log.Debugf("found %d github_organization_workflow_state records", len(t))
	for _, r := range t {
		ch <- prometheus.MustNewConstMetric(c.orgWorkflowState, prometheus.GaugeValue, 1.0, r.RepoOwner, r.RepoName, r.WorkflowName, r.WorkflowState)
	}

	limits, err := gh.GetRateLimits()
	if err != nil {
		log.Errorf("getting rate limit data for prometheus: %v", err)
	} else if limits.Core != nil {
		ch <- prometheus.MustNewConstMetric(c.rateLimitLimit, prometheus.GaugeValue, float64(limits.Core.Limit))
		ch <- prometheus.MustNewConstMetric(c.rateLimitUsed, prometheus.GaugeValue, float64(limits.Core.Limit-limits.Core.Remaining))
	}
}
