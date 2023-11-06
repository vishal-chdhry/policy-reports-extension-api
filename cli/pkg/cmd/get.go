package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/vishal-chdhry/policy-reports-extension-api/cli/pkg/utils"
	"github.com/vishal-chdhry/policy-reports-extension-api/client/pkg/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	watcher "k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/rest"
)

var watch bool

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "list resources",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("not enough arguments to 'get'")
		}
		if len(args) > 1 {
			return fmt.Errorf("too many arguments to 'get'")
		}
		d := newDoer(args[0])
		rv, err := d.getList(args[0])
		if err != nil {
			return err
		}
		if watch {
			return d.watch(args[0], rv)
		}
		return nil
	},
}

func init() {
	getCmd.Flags().BoolVarP(&watch, "watch", "w", false, "watch for changes")
}

type doer struct {
	client  *v1alpha1.DemoPolicyV1alpha2Client
	printer printers.ResourcePrinter
}

func newDoer(resource string) doer {
	printOpts := printers.PrintOptions{WithNamespace: true}
	return doer{
		client:  client,
		printer: printers.NewTypeSetter(v1alpha1.Scheme).ToPrinter(printers.NewTablePrinter(printOpts)),
	}
}

func (d doer) getList(resource string) (string, error) {
	var unst *unstructured.UnstructuredList
	if utils.IsClusterPolicyReport(resource) {
		cpol, err := d.client.ClusterPolicyReports().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return "", err
		}
		unst = utils.ClusterPolicyReportListToUnstructuredList(cpol)
	} else if utils.IsPolicyReport(resource) {
		pol, err := d.client.PolicyReports(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return "", err
		}
		unst = utils.PolicyReportListToUnstructuredList(pol)
	} else {
		return "", errors.New("unsupported resource")
	}

	table, err := toTable(unst.UnstructuredContent())
	if err != nil {
		return "", err
	}
	rv := unst.GetResourceVersion()
	if err = d.printer.PrintObj(table, os.Stdout); err != nil {
		return rv, err
	}
	// fmt.Fprint(os.Stdout, table)
	return rv, nil
}

func (d doer) watch(resource, rv string) error {
	var watcher watcher.Interface
	var err error
	if rv == "" {
		rv = "0"
	}
	opts := metav1.ListOptions{
		ResourceVersion: rv,
	}

	if utils.IsClusterPolicyReport(resource) {
		watcher, err = d.client.ClusterPolicyReports().Watch(context.TODO(), opts)
		if err != nil {
			return err
		}
	} else if utils.IsPolicyReport(resource) {
		watcher, err = d.client.PolicyReports(namespace).Watch(context.TODO(), opts)
		if err != nil {
			return err
		}
	} else {
		return errors.New("unsupported resource")
	}

	for event := range watcher.ResultChan() {
		table, err := toTable(event.Object.(*unstructured.Unstructured).UnstructuredContent())
		if err != nil {
			return err
		}
		if err = d.printer.PrintObj(table, os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

func toTable(obj map[string]interface{}) (*metav1.Table, error) {
	table := &metav1.Table{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj, table)
	if err != nil {
		return nil, err
	}
	for i := range table.Rows {
		row := &table.Rows[i]
		if row.Object.Raw == nil || row.Object.Object != nil {
			continue
		}
		converted, err := runtime.Decode(unstructured.UnstructuredJSONScheme, row.Object.Raw)
		if err != nil {
			return nil, err
		}
		row.Object.Object = converted
	}
	return table, nil
}

// endpointRestrictions implements negotiation.EndpointRestrictions.
type endpointRestrictions struct{}

// AllowsMediaTypeTransform implements negotiation.EndpointRestrictions.AllowsMediaTypeTransform.
func (e endpointRestrictions) AllowsMediaTypeTransform(mimeType string, mimeSubType string, gvk *schema.GroupVersionKind) bool {
	if mimeType != "application" {
		return false
	}
	if mimeSubType != "json" {
		return false
	}
	if gvk == nil || gvk.Kind == "" || gvk.Kind == "Table" {
		return true
	}
	return false
}

// AllowsServerVersion implements negotiation.EndpointRestrictions.AllowsServerVersion.
func (e endpointRestrictions) AllowsServerVersion(string) bool {
	return true
}

// AllowsStreamSchema implements negotiation.EndpointRestrictions.AllowsStreamSchema.
func (e endpointRestrictions) AllowsStreamSchema(string) bool {
	return true
}

type addOptions struct {
	accept string
	query  map[string]string
	next   http.RoundTripper
}

// RoundTrip adds parameters to request table formatting to the request object.
func (a *addOptions) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Accept", a.accept)
	q := r.URL.Query()
	for k, v := range a.query {
		q.Set(k, v)
	}
	r.URL.RawQuery = q.Encode()
	return a.next.RoundTrip(r)
}

func roundTripper(mediaType negotiation.MediaTypeOptions) func(http.RoundTripper) http.RoundTripper {
	accept := mediaType.Accepted.MediaType
	if mediaType.Convert != nil {
		if mediaType.Convert.Kind != "" {
			accept += ";as=" + mediaType.Convert.Kind
		}
		if mediaType.Convert.Version != "" {
			accept += ";v=" + mediaType.Convert.Version
		}
		if mediaType.Convert.Group != "" {
			accept += ";g=" + mediaType.Convert.Group
		}
	}
	return func(rt http.RoundTripper) http.RoundTripper {
		ao := addOptions{
			accept: accept,
			next:   rt,
		}
		if mediaType.Convert != nil && mediaType.Convert.Kind == "Table" {
			ao.query = map[string]string{
				"includeObject": "Object",
			}
		}
		return &ao
	}
}

func setTableMediaType(cfg *rest.Config) {
	acceptedTypes := []runtime.SerializerInfo{
		{
			MediaType:        "application/json",
			MediaTypeType:    "application",
			MediaTypeSubType: "json",
		},
		{
			MediaType:        "application/yaml",
			MediaTypeType:    "application",
			MediaTypeSubType: "yaml",
		},
	}
	mediaType, _ := negotiation.NegotiateMediaTypeOptions("application/json;as=Table;v=v1;g=meta.k8s.io", acceptedTypes, endpointRestrictions{})
	setOptions := roundTripper(mediaType)
	cfg.Wrap(setOptions)
}
