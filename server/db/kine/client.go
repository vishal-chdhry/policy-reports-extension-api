package kine

import (
	"github.com/k3s-io/kine/pkg/client"
	"github.com/k3s-io/kine/pkg/endpoint"
	"github.com/k3s-io/kine/pkg/tls"
	"github.com/vishal-chdhry/policy-reports-extension-api/server/pkg/common"
)

type clientConfigOpts func(*endpoint.ETCDConfig)

func buildKineOpts(opts ...clientConfigOpts) endpoint.ETCDConfig {
	cfg := endpoint.ETCDConfig{
		Endpoints: common.Endpoints,
		TLSConfig: tls.Config{
			CAFile:   common.LookupEnvOrDefault(common.CAEnvVar, common.CAFile),
			CertFile: common.LookupEnvOrDefault(common.CertEnvVar, common.CertFile),
			KeyFile:  common.LookupEnvOrDefault(common.KeyEnvVar, common.KeyFile),
		},
		LeaderElect: false,
	}

	for _, o := range opts {
		o(&cfg)
	}

	return cfg
}

func New(opts ...clientConfigOpts) (client.Client, error) {
	kClient, err := client.New(buildKineOpts(opts...))
	if err != nil {
		return nil, err
	}

	return kClient, nil
}

func WithEndpoints(endpoints []string) clientConfigOpts {
	return func(e *endpoint.ETCDConfig) {
		if len(endpoints) > 0 {
			e.Endpoints = endpoints
		}
	}
}

func WithKeyFile(file string) clientConfigOpts {
	return func(e *endpoint.ETCDConfig) {
		if len(file) > 0 {
			e.TLSConfig.KeyFile = file
		}
	}
}

func WithCAFile(file string) clientConfigOpts {
	return func(e *endpoint.ETCDConfig) {
		if len(file) > 0 {
			e.TLSConfig.CAFile = file
		}
	}
}

func WithCertFile(file string) clientConfigOpts {
	return func(e *endpoint.ETCDConfig) {
		if len(file) > 0 {
			e.TLSConfig.CertFile = file
		}
	}
}

func WithLeaderElection(val bool) clientConfigOpts {
	return func(e *endpoint.ETCDConfig) {
		e.LeaderElect = val
	}
}
