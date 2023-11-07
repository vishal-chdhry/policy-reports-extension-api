package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vishal-chdhry/policy-reports-extension-api/cli/pkg/utils"
	"github.com/vishal-chdhry/policy-reports-extension-api/client/pkg/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/printers"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply resources",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("too many arguments to 'apply'")
		}
		if len(inputfilepath) == 0 {
			return fmt.Errorf("input file path is required to 'create'")
		}
		d := newApplyDoer(args[0])
		var err error
		err = d.apply()
		if err != nil {
			return err
		}
		return nil
	},
}

type applyDoer struct {
	client  *v1alpha1.DemoPolicyV1alpha2Client
	printer printers.ResourcePrinter
}

func newApplyDoer(resource string) applyDoer {
	printOpts := printers.PrintOptions{WithNamespace: true}
	return applyDoer{
		client:  client,
		printer: printers.NewTypeSetter(v1alpha1.Scheme).ToPrinter(printers.NewTablePrinter(printOpts)),
	}
}

func (d applyDoer) apply() error {
	var unst *unstructured.Unstructured
	b, err := os.ReadFile(inputfilepath)
	if err != nil {
		return err
	}

	switch filepath.Ext(inputfilepath) {
	case ".yaml":
		unst = utils.YAMLToUnstructured(b)
	case ".json":
		unst = utils.JSONToUnstructured(b)
	default:
		return errors.New("unsupported file type")
	}

	polb, err := json.Marshal(unst.UnstructuredContent())
	if err != nil {
		return err
	}
	if unst.GetKind() == "ClusterPolicyReport" {
		var clusterPolicyReport *v1alpha1.ClusterPolicyReport
		err := json.Unmarshal(polb, clusterPolicyReport)
		if err != nil {
			return err
		}
		_, err = d.client.ClusterPolicyReports().Update(context.TODO(), clusterPolicyReport, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	} else if unst.GetKind() == "PolicyReport" {
		var policyReport *v1alpha1.PolicyReport
		err := json.Unmarshal(polb, policyReport)
		if err != nil {
			return err
		}
		_, err = d.client.PolicyReports(namespace).Update(context.TODO(), policyReport, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	} else {
		return errors.New("unsupported resource")
	}

	if len(unst.GetNamespace()) > 0 {
		fmt.Fprintln(os.Stdout, unst.GetKind(), " '", unst.GetName(), "' in namespace ", unst.GetNamespace(), " successfully configured.")
	} else {
		fmt.Fprintln(os.Stdout, unst.GetKind(), " '", unst.GetName(), "' successfully configured.")
	}
	return nil
}
