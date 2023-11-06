package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/k3s-io/kine/pkg/client"
	"github.com/vishal-chdhry/policy-reports-extension-api/client/pkg/v1alpha1"
	"github.com/vishal-chdhry/policy-reports-extension-api/server/pkg/common"
	"k8s.io/apimachinery/pkg/api/errors"
)

type ClusterPolicyReportsInterface interface {
	Get(ctx context.Context, name string) (*v1alpha1.ClusterPolicyReport, error)

	List(ctx context.Context) (*v1alpha1.ClusterPolicyReportList, error)

	Delete(ctx context.Context, name string) error

	DeleteCollection(ctx context.Context) error

	Update(ctx context.Context, clusterPolicyReport *v1alpha1.ClusterPolicyReport, name string) (*v1alpha1.ClusterPolicyReport, error)

	Create(ctx context.Context, clusterPolicyReport *v1alpha1.ClusterPolicyReport) (*v1alpha1.ClusterPolicyReport, error)
}

type clusterpolicyreportshandler struct {
	dbClient client.Client
}

func newClusterPolicyHandler(dbClient client.Client) ClusterPolicyReportsInterface {
	return &clusterpolicyreportshandler{
		dbClient: dbClient,
	}
}

func (c *clusterpolicyreportshandler) Get(ctx context.Context, name string) (*v1alpha1.ClusterPolicyReport, error) {
	if len(name) == 0 {
		return nil, errors.NewBadRequest("name  cannot be nil")
	}

	val, err := c.dbClient.Get(ctx, getClusterPolicyReportKey(name))
	if err != nil {
		return nil, err
	}

	var clusterPolicyReport v1alpha1.ClusterPolicyReport
	err = json.Unmarshal(val.Data, &clusterPolicyReport)
	if err != nil {
		return nil, errors.NewBadRequest("invalid object found")
	}
	return &clusterPolicyReport, nil
}

func (c *clusterpolicyreportshandler) List(ctx context.Context) (*v1alpha1.ClusterPolicyReportList, error) {
	val, err := c.dbClient.List(ctx, getClusterPolicyReportKeyForList(), 0) // TODO: Revision?
	if err != nil {
		return nil, err
	}

	var clusterPolicyReportList *v1alpha1.ClusterPolicyReportList = &v1alpha1.ClusterPolicyReportList{}
	clusterPolicyReportList.Kind = "ClusterPolicyReportList"
	clusterPolicyReportList.APIVersion = "prext.demo/v1alpha1"
	clusterPolicyReportList.ResourceVersion = fmt.Sprint(time.Now().Unix() % 100000)
	var clusterPolicyReports = make([]v1alpha1.ClusterPolicyReport, 0)
	for _, v := range val {
		var clusterPolicyReport v1alpha1.ClusterPolicyReport
		err = json.Unmarshal(v.Data, &clusterPolicyReport)
		if err != nil {
			return nil, errors.NewBadRequest("invalid object found")
		}
		clusterPolicyReports = append(clusterPolicyReports, clusterPolicyReport)
	}
	clusterPolicyReportList.Items = clusterPolicyReports

	return clusterPolicyReportList, nil
}

func (c *clusterpolicyreportshandler) Create(ctx context.Context, clusterPolicyReport *v1alpha1.ClusterPolicyReport) (*v1alpha1.ClusterPolicyReport, error) {
	if clusterPolicyReport == nil {
		return nil, errors.NewBadRequest("clusterpolicyreport cannot be nil")
	}

	b, err := json.Marshal(*clusterPolicyReport)
	if err != nil {
		return nil, errors.NewBadRequest("clusterpolicyreport could not be marshalled")
	}

	err = c.dbClient.Create(ctx, getClusterPolicyReportKey(clusterPolicyReport.Name), b)
	if err != nil {
		return nil, err
	}

	return clusterPolicyReport, nil
}

func (c *clusterpolicyreportshandler) Update(ctx context.Context, clusterPolicyReport *v1alpha1.ClusterPolicyReport, name string) (*v1alpha1.ClusterPolicyReport, error) {
	if clusterPolicyReport == nil || len(name) == 0 {
		return nil, errors.NewBadRequest("clusterpolicyreport or name cannot be nil")
	}

	b, err := json.Marshal(*clusterPolicyReport)
	if err != nil {
		return nil, errors.NewBadRequest("clusterpolicyreport could not be marshalled")
	}

	err = c.dbClient.Update(ctx, getClusterPolicyReportKey(name), 0, b) // TODO: Revision?
	if err != nil {
		return nil, err
	}

	return clusterPolicyReport, nil
}

func (c *clusterpolicyreportshandler) Delete(ctx context.Context, name string) error {
	if len(name) == 0 {
		return errors.NewBadRequest("name cannot be nil")
	}

	return c.dbClient.Delete(ctx, getClusterPolicyReportKey(name), 1) // TODO: Revision?
}

func (c *clusterpolicyreportshandler) DeleteCollection(ctx context.Context) error {
	list, err := c.List(ctx)
	if err != nil {
		return err
	}

	for _, v := range list.Items {
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
