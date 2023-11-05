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

type ClusterPolicyReportsInterface interface {
	Get(ctx context.Context, name string) (*v1alpha2.ClusterPolicyReport, error)

	List(ctx context.Context) ([]*v1alpha2.ClusterPolicyReport, error)

	Delete(ctx context.Context, name string) error

	DeleteCollection(ctx context.Context) error

	Update(ctx context.Context, clusterPolicyReport *v1alpha2.ClusterPolicyReport, name string) (*v1alpha2.ClusterPolicyReport, error)

	Create(ctx context.Context, clusterPolicyReport *v1alpha2.ClusterPolicyReport) (*v1alpha2.ClusterPolicyReport, error)
}

type clusterpolicyreportshandler struct {
	kineClient client.Client
}

func newClusterPolicyHandler(dbClient client.Client) ClusterPolicyReportsInterface {
	return &clusterpolicyreportshandler{
		kineClient: dbClient,
	}
}

func (c *clusterpolicyreportshandler) Get(ctx context.Context, name string) (*v1alpha2.ClusterPolicyReport, error) {
	if len(name) == 0 {
		return nil, errors.NewBadRequest("name  cannot be nil")
	}

	val, err := c.kineClient.Get(ctx, getClusterPolicyReportKey(name))
	if err != nil {
		return nil, err
	}

	var clusterPolicyReport v1alpha2.ClusterPolicyReport
	err = json.Unmarshal(val.Data, &clusterPolicyReport)
	if err != nil {
		return nil, errors.NewBadRequest("invalid object found")
	}
	return &clusterPolicyReport, nil
}

func (c *clusterpolicyreportshandler) List(ctx context.Context) ([]*v1alpha2.ClusterPolicyReport, error) {
	val, err := c.kineClient.List(ctx, getClusterPolicyReportKeyForList(), 0) // TODO: Revision?
	if err != nil {
		return nil, err
	}

	var clusterPolicyReports = make([]*v1alpha2.ClusterPolicyReport, 0)
	for _, v := range val {
		var clusterPolicyReport v1alpha2.ClusterPolicyReport
		err = json.Unmarshal(v.Data, &clusterPolicyReport)
		if err != nil {
			return nil, errors.NewBadRequest("invalid object found")
		}
		clusterPolicyReports = append(clusterPolicyReports, &clusterPolicyReport)
	}

	return clusterPolicyReports, nil
}

func (c *clusterpolicyreportshandler) Create(ctx context.Context, clusterPolicyReport *v1alpha2.ClusterPolicyReport) (*v1alpha2.ClusterPolicyReport, error) {
	if clusterPolicyReport == nil {
		return nil, errors.NewBadRequest("clusterpolicyreport cannot be nil")
	}

	b, err := json.Marshal(*clusterPolicyReport)
	if err != nil {
		return nil, errors.NewBadRequest("clusterpolicyreport could not be marshalled")
	}

	err = c.kineClient.Create(ctx, getClusterPolicyReportKey(clusterPolicyReport.Name), b)
	if err != nil {
		return nil, err
	}

	return clusterPolicyReport, nil
}

func (c *clusterpolicyreportshandler) Update(ctx context.Context, clusterPolicyReport *v1alpha2.ClusterPolicyReport, name string) (*v1alpha2.ClusterPolicyReport, error) {
	if clusterPolicyReport == nil || len(name) == 0 {
		return nil, errors.NewBadRequest("clusterpolicyreport or name cannot be nil")
	}

	b, err := json.Marshal(*clusterPolicyReport)
	if err != nil {
		return nil, errors.NewBadRequest("clusterpolicyreport could not be marshalled")
	}

	err = c.kineClient.Update(ctx, getClusterPolicyReportKey(name), 0, b) // TODO: Revision?
	if err != nil {
		return nil, err
	}

	return clusterPolicyReport, nil
}

func (c *clusterpolicyreportshandler) Delete(ctx context.Context, name string) error {
	if len(name) == 0 {
		return errors.NewBadRequest("name cannot be nil")
	}

	return c.kineClient.Delete(ctx, getClusterPolicyReportKey(name), 1) // TODO: Revision?
}

func (c *clusterpolicyreportshandler) DeleteCollection(ctx context.Context) error {
	list, err := c.List(ctx)
	if err != nil {
		return err
	}

	for _, v := range list {
		err := c.Delete(ctx, v.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func getClusterPolicyReportKey(name string) string {
	return fmt.Sprintf("/apis/%s/clusterpolicyreports/%s", common.GroupVersion, name)
}

func getClusterPolicyReportKeyForList() string {
	return fmt.Sprintf("/apis/%s/clusterpolicyreports/", common.GroupVersion)
}
