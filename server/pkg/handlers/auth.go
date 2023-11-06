package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	corev1listers "k8s.io/client-go/listers/core/v1"
)

const (
	FieldSelectorKey    = "fieldSelector"
	KubeSystemNamespace = "kube-system"
	ExtensionConfigMap  = "extension-apiserver-authentication"
	ClientCAKey         = "requestheader-client-ca-file"
	AllowedCNKey        = "requestheader-allowed-names"
)

func AuthenticateMiddleware(configMapCache corev1listers.ConfigMapNamespaceLister, logger *zap.SugaredLogger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(r.TLS.PeerCertificates) == 0 {
				logger.Warnf("user is not authenticated")
				http.Error(w, "user is not authenticated", http.StatusUnauthorized)
				return
			}
			requestCN := r.TLS.PeerCertificates[0].Subject.CommonName
			logger.Infof("authenticating user %s", requestCN)
			config, err := configMapCache.Get(ExtensionConfigMap)
			if err != nil {
				logger.Errorf("could not authenticate API server, err: %v", err)
				http.Error(w, fmt.Sprintf("could not authenticate API server, error: %v", err), http.StatusInternalServerError)
				return
			}
			allowedCNString, ok := config.Data[AllowedCNKey]
			if !ok {
				http.Error(w, "could not authenticate API server, invalid extension config", http.StatusInternalServerError)
			}
			allowedCN := []string{}
			if err := json.Unmarshal([]byte(allowedCNString), &allowedCN); err != nil {
				logger.Errorf("could not authenticate API server, err: %v", err)
				http.Error(w, fmt.Sprintf("could not authenticate API server, error: %v", err), http.StatusInternalServerError)
				return
			}
			found := false
			for _, allowed := range allowedCN {
				if allowed == requestCN {
					found = true
					break
				}
			}
			if !found {
				logger.Warnf("could not find user %s in allowed users", requestCN)
				http.Error(w, fmt.Sprintf("user %s not allowed", requestCN), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetClientCA(configMapCache corev1listers.ConfigMapLister) (string, error) {
	config, err := configMapCache.ConfigMaps(KubeSystemNamespace).Get(ExtensionConfigMap)
	if err != nil {
		return "", err
	}
	clientCA, ok := config.Data["requestheader-client-ca-file"]
	if !ok {
		return "", fmt.Errorf("invalid extension config")
	}
	return string(clientCA), nil
}
