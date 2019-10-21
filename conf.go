package main

import "os"

var (
	masterPassword        string
	hydraURL              string
	accountDirectory      string
	pathToStaticResources string
	pathToLoginPage       string
	pathToLoginFailedPage string
	port                  string
)

func envDefault(name string, backup string) string {
	found, exists := os.LookupEnv(name)
	if exists {
		return found
	}
	return backup
}

func init() {
	masterPassword = envDefault("CREAMY_MASTER_PASSWORD", "pancakes")
	hydraURL = envDefault("CREAMY_HYDRA_URL", "http://hydra.localhost:4445")
	accountDirectory = envDefault("CREAMY_ACCOUNTS_DIRECTORY", "./accounts")
	pathToStaticResources = envDefault("CREAMY_PATH_TO_STATIC_RESOURCES", "./sample-login/static")
	pathToLoginPage = envDefault("CREAMY_PATH_TO_LOGIN_PAGE", "./sample-login/login.html")
	pathToLoginFailedPage = envDefault("CREAMY_PATH_TO_LOGIN_FAILED_PAGE", "./sample-login/failed.html")
	port = envDefault("CREAMY_HTTP_PORT", "7000")
}