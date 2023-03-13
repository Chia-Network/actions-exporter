package prometheus

var (
	githubOrganization string
)

// Init accepts arguments relevant to the configuration of this package
func Init(org string) {
	githubOrganization = org
}
