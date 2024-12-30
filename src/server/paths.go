package server

const (
	ApiBase                       = "/api"
	ApiLogin                      = ApiBase + "/login"
	ApiLogout                     = ApiBase + "/logout"
	ApiRegister                   = ApiBase + "/register"
	ApiWinesDashboard             = ApiBase + "/wines/basic"
	ApiWinesWineCreation          = ApiBase + "/wines/create"
	ApiWinesFetchRegionsCountries = ApiBase + "/wines/countries/regions"
	ApiWinesFetchTypes            = ApiBase + "/wines/types"
	ApiWinesFetchBottleSizes      = ApiBase + "/wines/bottles/sizes"
	ApiWinesFetchDomains          = ApiBase + "/wines/domains"
	ApiWinesFetchPaginated        = ApiBase + "/wines"

	Frontend    = "/v1/*filepath"
	Favicon     = "/favicon.ico"
	RealFavicon = "/v1/icon/favicon.ico"
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
		ApiRegister,
	}
)

func init() {
	AuthProtectedPages = append(AuthProtectedPages, FrontendProtectedPages...)
}
