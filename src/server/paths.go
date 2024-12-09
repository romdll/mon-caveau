package server

const (
	ApiBase           = "/api"
	ApiLogin          = ApiBase + "/login"
	ApiWinesDashboard = ApiBase + "/wines/basic"

	Frontend = "/v1/*filepath"
)

var (
	FrontendProtectedPages = []string{
		"/dashboard",
	}

	AuthProtectedPages = []string{
		ApiBase,
	}
	AuthAvoidPages = []string{
		ApiLogin,
	}
)

func init() {
	AuthProtectedPages = append(AuthProtectedPages, FrontendProtectedPages...)
}
