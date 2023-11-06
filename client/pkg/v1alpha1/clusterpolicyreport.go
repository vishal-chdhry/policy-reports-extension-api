package v1alpha1

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ClusterPolicyReportsGetter has a method to return a ClusterPolicyReportInterface.
// A group's client should implement this interface.
type ClusterPolicyReportsGetter interface {
	ClusterPolicyReports() ClusterPolicyReportInterface
}

// ClusterPolicyReportInterface has methods to work with ClusterPolicyReport resources.
type ClusterPolicyReportInterface interface {
	Create(ctx context.Context, clusterPolicyReport *ClusterPolicyReport, opts v1.CreateOptions) (*ClusterPolicyReport, error)
	Update(ctx context.Context, clusterPolicyReport *ClusterPolicyReport, opts v1.UpdateOptions) (*ClusterPolicyReport, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*ClusterPolicyReport, error)
	List(ctx context.Context, opts v1.ListOptions) (*ClusterPolicyReportList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *ClusterPolicyReport, err error)
}

// clusterPolicyReports implements ClusterPolicyReportInterface
type clusterPolicyReports struct {
	client rest.Interface
}

// newClusterPolicyReports returns a ClusterPolicyReports
func newClusterPolicyReports(c *DemoPolicyV1alpha2Client) *clusterPolicyReports {
	return &clusterPolicyReports{
		client: c.RESTClient(),
	}
}

// Get takes name of the clusterPolicyReport, and returns the corresponding clusterPolicyReport object, and an error if there is any.
func (c *clusterPolicyReports) Get(ctx context.Context, name string, opts v1.GetOptions) (result *ClusterPolicyReport, err error) {
	result = &ClusterPolicyReport{}
	err = c.client.Get().
		Resource("clusterpolicyreports").
		Name(name).
		SpecificallyVersionedParams(&opts, ParameterCodec, GroupVersion).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterPolicyReports that match those selectors.
func (c *clusterPolicyReports) List(ctx context.Context, opts v1.ListOptions) (result *ClusterPolicyReportList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &ClusterPolicyReportList{}
	err = c.client.Get().
		Resource("clusterpolicyreports").
		SpecificallyVersionedParams(&opts, ParameterCodec, GroupVersion).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested clusterPolicyReports.
func (c *clusterPolicyReports) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("clusterpolicyreports").
		SpecificallyVersionedParams(&opts, ParameterCodec, GroupVersion).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a clusterPolicyReport and creates it.  Returns the server's representation of the clusterPolicyReport, and an error, if there is any.
func (c *clusterPolicyReports) Create(ctx context.Context, clusterPolicyReport *ClusterPolicyReport, opts v1.CreateOptions) (result *ClusterPolicyReport, err error) {
	result = &ClusterPolicyReport{}
	err = c.client.Post().
		Resource("clusterpolicyreports").
		SpecificallyVersionedParams(&opts, ParameterCodec, GroupVersion).
		Body(clusterPolicyReport).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a clusterPolicyReport and updates it. Returns the server's representation of the clusterPolicyReport, and an error, if there is any.
func (c *clusterPolicyReports) Update(ctx context.Context, clusterPolicyReport *ClusterPolicyReport, opts v1.UpdateOptions) (result *ClusterPolicyReport, err error) {
	result = &ClusterPolicyReport{}
	err = c.client.Put().
		Resource("clusterpolicyreports").
		Name(clusterPolicyReport.Name).
		SpecificallyVersionedParams(&opts, ParameterCodec, GroupVersion).
		Body(clusterPolicyReport).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the clusterPolicyReport and deletes it. Returns an error if one occurs.
func (c *clusterPolicyReports) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("clusterpolicyreports").
		SpecificallyVersionedParams(&opts, ParameterCodec, GroupVersion).
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *clusterPolicyReports) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("clusterpolicyreports").
		SpecificallyVersionedParams(&listOpts, ParameterCodec, GroupVersion).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched clusterPolicyReport.
func (c *clusterPolicyReports) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *ClusterPolicyReport, err error) {
	result = &ClusterPolicyReport{}
	err = c.client.Patch(pt).
		Resource("clusterpolicyreports").
		Name(name).
		SubResource(subresources...).
		SpecificallyVersionedParams(&opts, ParameterCodec, GroupVersion).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
