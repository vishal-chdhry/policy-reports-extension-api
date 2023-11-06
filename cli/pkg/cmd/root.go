package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vishal-chdhry/policy-reports-extension-api/client/pkg/v1alpha1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/scheme"
	apischeme "sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	rootCmd       *cobra.Command
	namespace     string
	output        string
	inputfilepath string
	GroupVersion  = schema.GroupVersion{Group: "prext.demo", Version: "v1alpha1"}
	schemeBuilder = &apischeme.Builder{GroupVersion: GroupVersion}
	mapper        meta.RESTMapper
	client        *v1alpha1.DemoPolicyV1alpha2Client
)

func init() {
	utilruntime.Must(schemeBuilder.AddToScheme(scheme.Scheme))

	kubecfgFlags := genericclioptions.NewConfigFlags(false)

	rootCmd = &cobra.Command{
		Use:   "kubectl-prext",
		Short: "cli for policy reports and cluster policy reports using aggregation api",
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			config, err := kubecfgFlags.ToRESTConfig()
			if err != nil {
				return err
			}

			mapper, err = kubecfgFlags.ToRESTMapper()
			if err != nil {
				return err
			}

			client, err = v1alpha1.NewForConfig(config)
			if err != nil {
				return err
			}

			return nil
		},
		SilenceUsage: true,
	}
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "parent namespace to scope request to")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output type, accepeted values: yaml,json")
	rootCmd.PersistentFlags().StringVarP(&inputfilepath, "file", "f", "", "path of input file")

	rootCmd.AddCommand(getCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
