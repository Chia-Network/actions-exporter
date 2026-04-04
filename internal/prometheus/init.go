package prometheus

var (
	githubOrganization string
	githubUser         string
)

// Init accepts arguments relevant to the configuration of this package
func Init(org, user string) {
	githubOrganization = org
	githubUser = user
}
