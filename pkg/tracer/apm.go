package tracer

import (
	"os"

	"github.com/dathuynh1108/clean-arch-base/pkg/config"
)

func InitAPMEnv() {
	var (
		config = config.GetConfig().APMConfig
	)
	os.Setenv("ELASTIC_APM_SERVER_URL", config.ServerURL)
	os.Setenv("ELASTIC_APM_SECRET_TOKEN", config.SecretToken)
	os.Setenv("ELASTIC_APM_SERVICE_NAME", config.ServiceName)
	os.Setenv("ELASTIC_APM_ENVIRONMENT", config.Environment)
}
