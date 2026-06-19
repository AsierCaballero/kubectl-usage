package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/AsierCaballero/kubectl-usage/pkg/analyzer"
	"github.com/AsierCaballero/kubectl-usage/pkg/collector"
	"github.com/AsierCaballero/kubectl-usage/pkg/output"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "kubectl-usage",
	Short: "Analyze pod resource usage against requests and limits",
	Long: `A kubectl plugin that analyzes pod resource consumption and
compares it against configured requests and limits. Detects
over-provisioned resources and provides downsizing recommendations.`,
	RunE: run,
}

var opts struct {
	namespace     string
	deployment    string
	labelSelector string
	allNamespaces bool
	sortBy        string
	format        string
	timeout       int
	verbose       bool
}

func init() {
	rootCmd.Flags().StringVarP(&opts.namespace, "namespace", "n", "default", "Kubernetes namespace")
	rootCmd.Flags().StringVarP(&opts.deployment, "deployment", "d", "", "Filter by deployment name")
	rootCmd.Flags().StringVarP(&opts.labelSelector, "selector", "l", "", "Label selector (e.g. app=myapp)")
	rootCmd.Flags().BoolVarP(&opts.allNamespaces, "all-namespaces", "A", false, "Show all namespaces")
	rootCmd.Flags().StringVarP(&opts.sortBy, "sort", "s", "waste", "Sort by: waste, cpu, memory, name, namespace")
	rootCmd.Flags().StringVarP(&opts.format, "output", "o", "table", "Output format: table, json")
	rootCmd.Flags().IntVarP(&opts.timeout, "timeout", "t", 30, "Timeout in seconds")
	rootCmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(opts.timeout)*time.Second)
	defer cancel()

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return fmt.Errorf("load kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("create kubernetes client: %w", err)
	}

	metricsClientset, err := metricsclient.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("create metrics client: %w", err)
	}

	c := collector.New(clientset, metricsClientset)
	collectOpts := collector.Options{
		Namespace:     opts.namespace,
		LabelSelector: opts.labelSelector,
		Deployment:    opts.deployment,
		AllNamespaces: opts.allNamespaces,
	}

	pods, err := c.Collect(ctx, collectOpts)
	if err != nil {
		return fmt.Errorf("collect: %w", err)
	}

	if len(pods) == 0 {
		fmt.Println("No pods found matching the criteria.")
		return nil
	}

	a := analyzer.New()
	results := a.Analyze(pods, analyzer.DefaultOptions)

	if err := output.Print(results, output.Format(opts.format), opts.sortBy); err != nil {
		return fmt.Errorf("output: %w", err)
	}

	return nil
}
