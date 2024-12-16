package server

const (
	ApiBase              = "/api"
	ApiLogin             = ApiBase + "/login"
	ApiLogout            = ApiBase + "/logout"
	ApiWinesDashboard    = ApiBase + "/wines/basic"
	ApiWinesWineCreation = ApiBase + "/wines/create"

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
