package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/vishal-chdhry/policy-reports-extension-api/server/db/inmemory"
	"github.com/vishal-chdhry/policy-reports-extension-api/server/pkg/common"
	"github.com/vishal-chdhry/policy-reports-extension-api/server/pkg/handlers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var (
	Namespace             = os.Getenv("POD_NAMESPACE")
	PodName               = os.Getenv("POD_NAME")
	CertManagerSecretName = common.LookupEnvOrDefault("CERT_MANAGER_SECRET", "prext-cert-secret")
	ServiceName           = common.LookupEnvOrDefault("SERVICE_NAME", "prext-svc")
	DeploymentName        = common.LookupEnvOrDefault("DEPLOYMENT_NAME", "prext-server")
	CertPath              = common.LookupEnvOrDefault("CERTPATH", "/certs/tls.crt")
	KeyPath               = common.LookupEnvOrDefault("KEYPATH", "/certs/tls.key")

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
	flag.StringVar(&host, "host", "", "Host to run the service on")

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

	slog.Info("getting kuberntes cluster config")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("failed to get kubernetes cluster config: %v", err)
	}

	slog.Info("staring kubernetes client")
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to initialize kube client: %v", err)
	}

	cmStopCh := make(chan struct{}, 1)
	cmInformer := common.NewConfigMapInformer(kubeClient, handlers.KubeSystemNamespace, handlers.ExtensionConfigMap, resyncPeriod)
	go cmInformer.Informer().Run(cmStopCh)
	if !cache.WaitForCacheSync(cmStopCh, cmInformer.Informer().HasSynced) {
		log.Fatalf("config map cache failed to sync")
		return
	}

	tlsStopCh := make(chan struct{}, 1)
	tlsInformer := common.NewSecretInformer(kubeClient, Namespace, CertManagerSecretName, resyncPeriod)
	go tlsInformer.Informer().Run(tlsStopCh)
	if !cache.WaitForCacheSync(tlsStopCh, tlsInformer.Informer().HasSynced) {
		log.Fatalf("tls secret cache failed to sync")
		return
	}

	// kineClient, err := kine.New(
	// 	kine.WithCAFile(dbCAFile),
	// 	kine.WithCertFile(dbCertFile),
	// 	kine.WithKeyFile(dbKeyFile),
	// 	kine.WithEndpoints(strings.Split(endpoints, ",")),
	// )
	// if err != nil {
	// 	log.Fatalf("failed to initialize kineclient: %v", err)
	// }

	slog.Info("starting inmemory database")
	inMemoryDb := inmemory.New(slog)
	handlerSet := handlers.NewHandlerSet(inMemoryDb, slog)

	slog.Info("setting up routing")
	mux := mux.NewRouter()
	mux.HandleFunc(fmt.Sprintf("/apis/%s", common.GroupVersion), handlers.TestHandler)
	mux.HandleFunc(fmt.Sprintf("/apis/%s/{resource}", common.GroupVersion), handlerSet.ClusterScopedHandler(ctx))
	mux.HandleFunc(fmt.Sprintf("/apis/%s/{resource}/{name}", common.GroupVersion), handlerSet.ClusterScopedHandlerWithName(ctx))
	mux.HandleFunc(fmt.Sprintf("/apis/%s/namespaces/{namespace}/{resource}", common.GroupVersion), handlerSet.NamespacedHandler(ctx))
	mux.HandleFunc(fmt.Sprintf("/apis/%s/namespaces/{namespace}/{resource}/{name}", common.GroupVersion), handlerSet.NamespacedHandlerWithName(ctx))
	mux.Use(handlers.AuthenticateMiddleware(cmInformer.Lister().ConfigMaps(handlers.KubeSystemNamespace), slog))

	address := host + ":" + fmt.Sprint(port)
	clientCA, err := handlers.GetClientCA(cmInformer.Lister())
	if err != nil {
		slog.Fatalf("failed to get client ca: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(clientCA))
	slog.Infof("starting server on port %s", address)
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			secret, err := tlsInformer.Lister().Secrets(Namespace).Get(CertManagerSecretName)
			if err != nil {
				return nil, err
			} else if secret == nil {
				return nil, errors.New("tls secret not found")
			} else if secret.Type != corev1.SecretTypeTLS {
				return nil, errors.New("secret is not a TLS secret")
			}

			TLSCert, found := secret.Data[corev1.TLSCertKey]
			if !found {
				return nil, errors.New("secret does not have TLS Cert")
			}
			TLSKey, found := secret.Data[corev1.TLSPrivateKeyKey]
			if !found {
				return nil, errors.New("secret does not have TLS Key")
			}

			cert, err := tls.X509KeyPair(TLSCert, TLSKey)
			if err != nil {
				return nil, err
			}

			return &cert, nil
		},
	}
	server := http.Server{
		Addr:      address,
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	err = server.ListenAndServeTLS("", "")
	slog.Infof("server started on port %s", address)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
