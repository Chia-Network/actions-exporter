package prometheus

import (
	"github.com/chia-network/actions-exporter/internal/github"
	log "github.com/sirupsen/logrus"
)

type workflowStateData struct {
	RepoOwner     string
	RepoName      string
	WorkflowName  string
	WorkflowState string
}

func getOrgWorkflowState() []workflowStateData {
	var r []workflowStateData

	repos, err := github.ListRepositoriesByOrg(githubOrganization)
	if err != nil {
		log.Errorf("getting workflow state data for prometheus: %v", err)
		return []workflowStateData{}
	}

	for _, repo := range repos {
		wfs, err := github.ListRepositoryWorkflows(*repo.Owner.Login, *repo.Name)
		if err != nil {
			log.Errorf("getting workflow state data for prometheus: %v", err)
		}

		for _, wf := range wfs {
			// Skip active workflows
			if *wf.State == "active" {
				continue
			}
			r = append(r, workflowStateData{
				RepoOwner:     *repo.Owner.Login,
				RepoName:      *repo.Name,
				WorkflowName:  *wf.Name,
				WorkflowState: *wf.State,
			})
		}
	}

	return r
}
