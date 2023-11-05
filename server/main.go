package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-logr/zapr"
	"github.com/gorilla/mux"
	"github.com/kyverno/pkg/certmanager"
	tlsMgr "github.com/kyverno/pkg/tls"
	"github.com/vishal-chdhry/policy-reports-extension-api/server/db/inmemory"
	"github.com/vishal-chdhry/policy-reports-extension-api/server/pkg/common"
	"github.com/vishal-chdhry/policy-reports-extension-api/server/pkg/handlers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	Namespace      = os.Getenv("POD_NAMESPACE")
	PodName        = os.Getenv("POD_NAME")
	ServiceName    = common.LookupEnvOrDefault("SERVICE_NAME", "svc")
	DeploymentName = common.LookupEnvOrDefault("DEPLOYMENT_NAME", "prext-server")

	CertRenewalInterval = 12 * time.Hour
	CAValidityDuration  = 365 * 24 * time.Hour
	TLSValidityDuration = 150 * 24 * time.Hour

	resyncPeriod = 15 * time.Minute
)

func main() {
	ctx := context.Background()

	var endpoints string
	flag.StringVar(&endpoints, "dbEndpoints", "", "Endpoints of the database.")

	var dbCAFile string
	flag.StringVar(&dbCAFile, "dbCAFile", "", "Ca file location of the database.")

	var dbCertFile string
	flag.StringVar(&dbCertFile, "dbCertFile", "", "Cert file location of the database.")

	var dbKeyFile string
	flag.StringVar(&dbKeyFile, "dbKeyFile", "", "Key file location of the database.")

	var host string
	flag.StringVar(&host, "host", "127.0.0.1", "Host to run the service on")

	var port int
	flag.IntVar(&port, "port", 7443, "Port to run the service on")

	flag.Parse()

	zc := zap.NewDevelopmentConfig()
	zc.Level = zap.NewAtomicLevelAt(zapcore.Level(-2))
	logger, err := zc.Build()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	logger = logger.WithOptions(zap.AddStacktrace(zap.DPanicLevel))
	slog := logger.Sugar()

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("failed to get kubernetes cluster config: %v", err)
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to initialize kube client: %v", err)
	}

	tlsMgrConfig := &tlsMgr.Config{
		ServiceName: ServiceName,
		Namespace:   Namespace,
	}

	caStopCh := make(chan struct{}, 1)
	caInformer := common.NewSecretInformer(kubeClient, Namespace, tlsMgr.GenerateRootCASecretName(tlsMgrConfig), resyncPeriod)
	go caInformer.Informer().Run(caStopCh)

	tlsStopCh := make(chan struct{}, 1)
	tlsInformer := common.NewSecretInformer(kubeClient, Namespace, tlsMgr.GenerateTLSPairSecretName(tlsMgrConfig), resyncPeriod)
	go tlsInformer.Informer().Run(tlsStopCh)

	certRenewer := tlsMgr.NewCertRenewer(
		zapr.NewLogger(logger).WithName("tls").WithValues("pod", PodName),
		kubeClient.CoreV1().Secrets(Namespace),
		CertRenewalInterval,
		CAValidityDuration,
		TLSValidityDuration,
		"",
		tlsMgrConfig,
	)

	certManager := certmanager.NewController(
		zapr.NewLogger(logger).WithName("certmanager").WithValues("pod", PodName),
		caInformer,
		tlsInformer,
		certRenewer,
		tlsMgrConfig,
	)

	go func() {
		certManager.Run(ctx, 1)
	}()

	// kineClient, err := kine.New(
	// 	kine.WithCAFile(dbCAFile),
	// 	kine.WithCertFile(dbCertFile),
	// 	kine.WithKeyFile(dbKeyFile),
	// 	kine.WithEndpoints(strings.Split(endpoints, ",")),
	// )
	// if err != nil {
	// 	log.Fatalf("failed to initialize kineclient: %v", err)
	// }

	inMemoryDb := inmemory.New(slog)
	handlerSet := handlers.NewHandlerSet(inMemoryDb)

	mux := mux.NewRouter()
	mux.HandleFunc(fmt.Sprintf("/apis/%s", common.GroupVersion), handlers.TestHandler)
	mux.HandleFunc(fmt.Sprintf("/apis/%s/{resource}", common.GroupVersion), handlerSet.ClusterScopedHandler(ctx))
	mux.HandleFunc(fmt.Sprintf("/apis/%s/{resource}/{name}", common.GroupVersion), handlerSet.ClusterScopedHandlerWithName(ctx))
	mux.HandleFunc(fmt.Sprintf("/apis/%s/namespaces/{namespace}/{resource}", common.GroupVersion), handlerSet.NamespacedHandler(ctx))
	mux.HandleFunc(fmt.Sprintf("/apis/%s/namespaces/{namespace}/{resource}/{name}", common.GroupVersion), handlerSet.NamespacedHandlerWithName(ctx))

	address := host + ":" + fmt.Sprint(port)
	tlsConf := &tls.Config{
		GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			secret, err := tlsInformer.Lister().Secrets(tlsMgrConfig.Namespace).Get(tlsMgr.GenerateTLSPairSecretName(tlsMgrConfig))
			if err != nil {
				return nil, err
			} else if secret == nil {
				return nil, errors.New("tls secret not found")
			} else if secret.Type != corev1.SecretTypeTLS {
				return nil, errors.New("secret is not a TLS secret")
			}

			cert, err := tls.X509KeyPair(secret.Data[corev1.TLSCertKey], secret.Data[corev1.TLSPrivateKeyKey])
			if err != nil {
				return nil, err
			}

			return &cert, nil
		},
	}
	server := http.Server{
		Addr:      address,
		Handler:   mux,
		TLSConfig: tlsConf,
	}

	fmt.Println("starting server on ", address)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
