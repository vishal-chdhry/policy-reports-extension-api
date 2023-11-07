package v1alpha1

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PolicyReportsGetter has a method to return a PolicyReportInterface.
// A group's client should implement this interface.
type PolicyReportsGetter interface {
	PolicyReports(namespace string) PolicyReportInterface
}

// PolicyReportInterface has methods to work with PolicyReport resources.
type PolicyReportInterface interface {
	Create(ctx context.Context, obj *unstructured.Unstructured, opts v1.CreateOptions) (*unstructured.Unstructured, error)
	Update(ctx context.Context, obj *unstructured.Unstructured, opts v1.UpdateOptions) (*unstructured.Unstructured, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*unstructured.Unstructured, error)
	List(ctx context.Context, opts v1.ListOptions) (*unstructured.UnstructuredList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *unstructured.Unstructured, err error)
}

// policyReports implements PolicyReportInterface
type policyReports struct {
	client rest.Interface
	ns     string
}

// newPolicyReports returns a PolicyReports
func newPolicyReports(c *DemoPolicyV1alpha2Client, namespace string) PolicyReportInterface {
	return &policyReports{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the policyReport, and returns the corresponding policyReport object, and an error if there is any.
func (c *policyReports) Get(ctx context.Context, name string, options v1.GetOptions) (result *unstructured.Unstructured, err error) {
	result = &unstructured.Unstructured{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("policyreports").
		Name(name).
		SpecificallyVersionedParams(&options, ParameterCodec, GroupVersion).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PolicyReports that match those selectors.
func (c *policyReports) List(ctx context.Context, opts v1.ListOptions) (result *unstructured.UnstructuredList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &unstructured.UnstructuredList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("policyreports").
		SpecificallyVersionedParams(&opts, ParameterCodec, GroupVersion).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested policyReports.
func (c *policyReports) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("policyreports").
		SpecificallyVersionedParams(&opts, ParameterCodec, GroupVersion).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a policyReport and creates it.  Returns the server's representation of the policyReport, and an error, if there is any.
func (c *policyReports) Create(ctx context.Context, obj *unstructured.Unstructured, opts v1.CreateOptions) (result *unstructured.Unstructured, err error) {
	result = &unstructured.Unstructured{}
	intr := obj.UnstructuredContent()
	body, err := json.Marshal(intr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get body")
	}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("policyreports").
		SpecificallyVersionedParams(&opts, ParameterCodec, GroupVersion).
		Body(string(body)).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a policyReport and updates it. Returns the server's representation of the policyReport, and an error, if there is any.
func (c *policyReports) Update(ctx context.Context, obj *unstructured.Unstructured, opts v1.UpdateOptions) (result *unstructured.Unstructured, err error) {
	result = &unstructured.Unstructured{}
	intr := obj.UnstructuredContent()
	body, err := json.Marshal(intr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get body")
	}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("policyreports").
		Name(obj.GetName()).
		SpecificallyVersionedParams(&opts, ParameterCodec, GroupVersion).
		Body(string(body)).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the policyReport and deletes it. Returns an error if one occurs.
func (c *policyReports) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("policyreports").
		Name(name).
		SpecificallyVersionedParams(&opts, ParameterCodec, GroupVersion).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *policyReports) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("policyreports").
		SpecificallyVersionedParams(&listOpts, ParameterCodec, GroupVersion).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched policyReport.
func (c *policyReports) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *unstructured.Unstructured, err error) {
	result = &unstructured.Unstructured{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("policyreports").
		Name(name).
		SubResource(subresources...).
		SpecificallyVersionedParams(&opts, ParameterCodec, GroupVersion).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
