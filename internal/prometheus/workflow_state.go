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

func getWorkflowState() []workflowStateData {
	var r []workflowStateData

	if githubOrganization != "" {
		r = append(r, getOrgWorkflowState()...)
	}
	if githubUser != "" {
		r = append(r, getUserWorkflowState()...)
	}

	return r
}

func getOrgWorkflowState() []workflowStateData {
	var r []workflowStateData

	repos, err := github.ListRepositoriesByOrg(githubOrganization)
	if err != nil {
		log.Errorf("getting org workflow state data for prometheus: %v", err)
		return nil
	}

	for _, repo := range repos {
		wfs, err := github.ListRepositoryWorkflows(*repo.Owner.Login, *repo.Name)
		if err != nil {
			log.Errorf("getting workflow state data for prometheus: %v", err)
			continue
		}

		for _, wf := range wfs {
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

func getUserWorkflowState() []workflowStateData {
	var r []workflowStateData

	repos, err := github.ListRepositoriesByUser(githubUser)
	if err != nil {
		log.Errorf("getting user workflow state data for prometheus: %v", err)
		return nil
	}

	for _, repo := range repos {
		wfs, err := github.ListRepositoryWorkflows(*repo.Owner.Login, *repo.Name)
		if err != nil {
			log.Errorf("getting workflow state data for prometheus: %v", err)
			continue
		}

		for _, wf := range wfs {
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
