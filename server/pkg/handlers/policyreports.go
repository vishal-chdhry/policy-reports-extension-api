package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/k3s-io/kine/pkg/client"
	"github.com/kyverno/kyverno/api/policyreport/v1alpha2"
	"github.com/vishal-chdhry/policy-reports-extension-api/server/pkg/common"
	"k8s.io/apimachinery/pkg/api/errors"
)

type PolicyReportsInterface interface {
	Get(ctx context.Context, name, namespace string) (*v1alpha2.PolicyReport, error)

	List(ctx context.Context, namespace string) ([]*v1alpha2.PolicyReport, error)

	Delete(ctx context.Context, name, namespace string) error

	DeleteCollection(ctx context.Context, namespace string) error

	Create(ctx context.Context, policyReport *v1alpha2.PolicyReport, namespace string) (*v1alpha2.PolicyReport, error)

	Update(ctx context.Context, policyReport *v1alpha2.PolicyReport, name, namespace string) (*v1alpha2.PolicyReport, error)
}

type policyreportshandler struct {
	kineClient client.Client
}

func newPolicyHandler(kineClient client.Client) PolicyReportsInterface {
	return &policyreportshandler{
		kineClient: kineClient,
	}
}

func (p *policyreportshandler) Get(ctx context.Context, name, namespace string) (*v1alpha2.PolicyReport, error) {
	if len(name) == 0 || len(namespace) == 0 {
		return nil, errors.NewBadRequest("name or namespace cannot be nil")
	}

	val, err := p.kineClient.Get(ctx, getPolicyReportKey(namespace, name))
	if err != nil {
		return nil, err
	}

	var policyReport v1alpha2.PolicyReport
	err = json.Unmarshal(val.Data, &policyReport)
	if err != nil {
		return nil, errors.NewBadRequest("invalid object found")
	}
	return &policyReport, nil
}

func (p *policyreportshandler) List(ctx context.Context, namespace string) ([]*v1alpha2.PolicyReport, error) {
	if len(namespace) == 0 {
		return nil, errors.NewBadRequest("namespace cannot be nil")
	}

	val, err := p.kineClient.List(ctx, getPolicyReportKeyForList(namespace), 0) // TODO: Revision?
	if err != nil {
		return nil, err
	}

	var policyReports = make([]*v1alpha2.PolicyReport, 0)
	for _, v := range val {
		var policyReport v1alpha2.PolicyReport
		err = json.Unmarshal(v.Data, &policyReport)
		if err != nil {
			return nil, errors.NewBadRequest("invalid object found")
		}
		policyReports = append(policyReports, &policyReport)
	}

	return policyReports, nil
}

func (p *policyreportshandler) Create(ctx context.Context, policyReport *v1alpha2.PolicyReport, namespace string) (*v1alpha2.PolicyReport, error) {
	if policyReport == nil || len(namespace) == 0 {
		return nil, errors.NewBadRequest("policyreport or namespace cannot be nil")
	}

	b, err := json.Marshal(*policyReport)
	if err != nil {
		return nil, errors.NewBadRequest("policyreport could not be marshalled")
	}

	err = p.kineClient.Create(ctx, getPolicyReportKey(namespace, policyReport.Name), b)
	if err != nil {
		return nil, err
	}

	return policyReport, nil
}

func (p *policyreportshandler) Update(ctx context.Context, policyReport *v1alpha2.PolicyReport, name, namespace string) (*v1alpha2.PolicyReport, error) {
	if policyReport == nil || len(name) == 0 || len(namespace) == 0 {
		return nil, errors.NewBadRequest("policyreport, name or namespace cannot be nil")
	}

	b, err := json.Marshal(*policyReport)
	if err != nil {
		return nil, errors.NewBadRequest("policyreport could not be marshalled")
	}

	err = p.kineClient.Update(ctx, getPolicyReportKey(namespace, name), 0, b) // TODO: Revision?
	if err != nil {
		return nil, err
	}

	return policyReport, nil
}

func (p *policyreportshandler) Delete(ctx context.Context, name, namespace string) error {
	if len(name) == 0 || len(namespace) == 0 {
		return errors.NewBadRequest("name or namespace cannot be nil")
	}

	return p.kineClient.Delete(ctx, getPolicyReportKey(namespace, name), 1) // TODO: Revision?
}

func (p *policyreportshandler) DeleteCollection(ctx context.Context, namespace string) error {
	if len(namespace) == 0 {
		return errors.NewBadRequest("namespace cannot be nil")
	}

	list, err := p.List(ctx, namespace)
	if err != nil {
		return err
	}

	for _, v := range list {
		err := p.Delete(ctx, v.Name, namespace)
		if err != nil {
			return err
		}
	}
	return nil
}

func getPolicyReportKey(namespace, name string) string {
	return fmt.Sprintf("/apis/%s/namespaces/%s/policyreports/%s", common.GroupVersion, namespace, name)
}

func getPolicyReportKeyForList(namespace string) string {
	return fmt.Sprintf("/apis/%s/namespaces/%s/policyreports/", common.GroupVersion, namespace)
}
