package common

const (
	Group        = "policyreports.vishal.demo"
	Version      = "v1alpha1"
	GroupVersion = Group + "/" + Version

	PORT     = 7443
	CAFile   = ""
	CertFile = ""
	KeyFile  = ""

	CAEnvVar   = "DB_CA_FILE"
	CertEnvVar = "DB_CERT_FILE"
	KeyEnvVar  = "DB_KEY_FILE"
)

var (
	Endpoints = make([]string, 0)
)
