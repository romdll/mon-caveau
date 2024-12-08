package server

const (
	ApiBase  = "/api"
	ApiLogin = ApiBase + "/login"

	Frontend = "*filepath"
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
