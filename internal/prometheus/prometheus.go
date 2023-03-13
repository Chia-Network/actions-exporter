package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Collector contains all prometheus metric Descs
type Collector struct {
	orgWorkflowState *prometheus.Desc
}

// NewCollector constructor function for Collector
func NewCollector() *Collector {
	return &Collector{
		orgWorkflowState: prometheus.NewDesc("github_organization_workflow_state",
			"Shows non-active workflow state for workflows in a GitHub organization.",
			[]string{"owner", "repository", "workflow", "state"}, nil,
		),
	}
}

// Describe contains all the prometheus descriptors for this metric collector
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.orgWorkflowState
}

// Collect instructs the prometheus client how to collect the metrics for each descriptor
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	t := getOrgWorkflowState()
	log.Debugf("found %d github_organization_workflow_state records", len(t))
	for _, r := range t {
		ch <- prometheus.MustNewConstMetric(c.orgWorkflowState, prometheus.GaugeValue, 1.0, r.RepoOwner, r.RepoName, r.WorkflowName, r.WorkflowState)
	}
}
