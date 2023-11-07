package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vishal-chdhry/policy-reports-extension-api/cli/pkg/utils"
	"github.com/vishal-chdhry/policy-reports-extension-api/client/pkg/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete resources",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("not enough arguments to 'delete'")
		}
		if len(args) > 2 {
			return fmt.Errorf("too many arguments to 'delete'")
		}
		d := newDeleteDoer(args[0])
		err := d.delete(args[0], args[1])
		if err != nil {
			return err
		}
		return nil
	},
}

type deleteDoer struct {
	client  *v1alpha1.DemoPolicyV1alpha2Client
	printer printers.ResourcePrinter
}

func newDeleteDoer(resource string) deleteDoer {
	printOpts := printers.PrintOptions{WithNamespace: true}
	return deleteDoer{
		client:  client,
		printer: printers.NewTypeSetter(v1alpha1.Scheme).ToPrinter(printers.NewTablePrinter(printOpts)),
	}
}

func (d deleteDoer) delete(resource, name string) error {
	if utils.IsClusterPolicyReport(resource) {
		err := d.client.ClusterPolicyReports().Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	} else if utils.IsPolicyReport(resource) {
		err := d.client.PolicyReports(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	} else {
		return errors.New("unsupported resource")
	}

	fmt.Fprintln(os.Stdout, "Successfully deleted ", resource, ": ", name)
	return nil
}
