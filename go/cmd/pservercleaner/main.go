package main

import (
	"context"
	"flag"

	log "github.com/inconshreveable/log15"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	logLevel := flag.String("log_level", "info", "Log level can be debug, info, warn, error, crit.")
	flag.Parse()

	lvl, err := log.LvlFromString(*logLevel)
	if err != nil {
		panic(err)
	}

	log.Root().SetHandler(
		log.LvlFilterHandler(lvl, log.CallerStackHandler("%+v", log.StderrHandler)),
	)

	// Create the client config. Use kubeconfig if given, otherwise assume in-cluster.
	config, err := buildConfig(*kubeconfig)
	if err != nil {
		panic(err)
	}

	// setup some optional configuration
	// paddlejob.ConfigureClient(config)

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	client, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	cleaner := NewCleaner(client, clientset)
	if err != nil {
		panic(err)
	}

	defer cancelFunc()
	cleaner.Run(ctx)
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
