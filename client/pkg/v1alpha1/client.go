package v1alpha1

import (
	"net/http"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/scheme"
)

type DemoPolicyV1alpha2Interface interface {
	RESTClient() rest.Interface
	PolicyReportsGetter
	ClusterPolicyReportsGetter
}

var (
	GroupVersion = schema.GroupVersion{Group: "prext.demo", Version: "v1alpha1"}
)

type DemoPolicyV1alpha2Client struct {
	restClient rest.Interface
}

func (c *DemoPolicyV1alpha2Client) ClusterPolicyReports() ClusterPolicyReportInterface {
	return newClusterPolicyReports(c)
}

func (c *DemoPolicyV1alpha2Client) PolicyReports(namespace string) PolicyReportInterface {
	return newPolicyReports(c, namespace)
}

func NewForConfig(c *rest.Config) (*DemoPolicyV1alpha2Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	httpClient, err := rest.HTTPClientFor(&config)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(&config, httpClient)
}

func NewForConfigAndClient(c *rest.Config, h *http.Client) (*DemoPolicyV1alpha2Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientForConfigAndClient(&config, h)
	if err != nil {
		return nil, err
	}
	return &DemoPolicyV1alpha2Client{client}, nil
}

// NewForConfigOrDie creates a new Wgpolicyk8sV1alpha2Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *DemoPolicyV1alpha2Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

func setConfigDefaults(config *rest.Config) error {
	gv := schema.GroupVersion{Group: "prext.demo", Version: "v1alpha1"}
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

func (c *DemoPolicyV1alpha2Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
