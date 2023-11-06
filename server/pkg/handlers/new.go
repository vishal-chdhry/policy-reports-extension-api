package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/k3s-io/kine/pkg/client"
	"github.com/vishal-chdhry/policy-reports-extension-api/client/pkg/v1alpha1"

	"go.uber.org/zap"
)

type HandlerSetInterface interface {
	ClusterScopedHandler(context.Context) func(http.ResponseWriter, *http.Request)
	ClusterScopedHandlerWithName(context.Context) func(http.ResponseWriter, *http.Request)
	NamespacedHandler(context.Context) func(http.ResponseWriter, *http.Request)
	NamespacedHandlerWithName(context.Context) func(http.ResponseWriter, *http.Request)
}

type handlerSet struct {
	logger      *zap.SugaredLogger
	cpolHandler ClusterPolicyReportsInterface
	polHandler  PolicyReportsInterface
}

func NewHandlerSet(dbClient client.Client, logger *zap.SugaredLogger) HandlerSetInterface {
	return &handlerSet{
		logger:      logger,
		cpolHandler: newClusterPolicyHandler(dbClient),
		polHandler:  newPolicyHandler(dbClient),
	}
}

func (h *handlerSet) ClusterScopedHandler(ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		h.logger.Infof("ClusterScopedHandler: recieved request method:%s, url:%s, resource:%s", r.Method, r.URL.String(), vars["resource"])

		if vars["resource"] != "clusterpolicyreports" {
			http.Error(w, "only clusterpolicyreports are supported", http.StatusBadRequest)
		}

		var data []byte = make([]byte, 0)
		var err error
		switch r.Method {
		case http.MethodGet:
			var val *v1alpha1.ClusterPolicyReportList
			val, err = h.cpolHandler.List(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotAcceptable)
			}

			data, err = json.MarshalIndent(val, "  ", "  ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodPost, http.MethodPut:
			var clusterPolicyReport *v1alpha1.ClusterPolicyReport
			raw, _ := io.ReadAll(r.Body)
			h.logger.Info("Body:", string(raw))
			err := json.Unmarshal(raw, &clusterPolicyReport)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotAcceptable)
			}
			clusterPolicyReport, err = h.cpolHandler.Create(ctx, clusterPolicyReport)

			data, err = json.MarshalIndent(clusterPolicyReport, "  ", "  ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodDelete:
			err = h.cpolHandler.DeleteCollection(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "method not supported on this endpoint", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func (h *handlerSet) ClusterScopedHandlerWithName(ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		h.logger.Infof("ClusterScopedHandlerWithName: recieved request method:%s, url:%s, resource:%s, name:%s", r.Method, r.URL.String(), vars["resource"], vars["name"])

		if vars["resource"] != "clusterpolicyreports" {
			http.Error(w, "only clusterpolicyreports are supported on this endpoint", http.StatusBadRequest)
		}

		var data []byte = make([]byte, 0)
		var err error
		switch r.Method {
		case http.MethodGet:
			var val *v1alpha1.ClusterPolicyReport
			val, err = h.cpolHandler.Get(ctx, vars["name"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotAcceptable)
			}

			data, err = json.MarshalIndent(val, "  ", "  ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodPost, http.MethodPut:
			var clusterPolicyReport *v1alpha1.ClusterPolicyReport
			raw, _ := io.ReadAll(r.Body)
			h.logger.Info("Body:", string(raw))
			err := json.Unmarshal(raw, &clusterPolicyReport)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotAcceptable)
			}
			clusterPolicyReport, err = h.cpolHandler.Update(ctx, clusterPolicyReport, vars["name"])

			data, err = json.MarshalIndent(clusterPolicyReport, "  ", "  ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodDelete:
			err = h.cpolHandler.Delete(ctx, vars["name"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "method not supported on this endpoint", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func (h *handlerSet) NamespacedHandler(ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		h.logger.Infof("NamespacedHandler: recieved request method:%s, url:%s, resource:%s, namespace:%s", r.Method, r.URL.String(), vars["resource"], vars["namespace"])

		if vars["resource"] != "policyreports" {
			http.Error(w, "only policyreports are supported on this endpoint", http.StatusBadRequest)
		}

		var data []byte = make([]byte, 0)
		var err error
		switch r.Method {
		case http.MethodGet:
			var val *v1alpha1.PolicyReportList
			val, err = h.polHandler.List(ctx, vars["namespace"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotAcceptable)
			}

			data, err = json.MarshalIndent(val, "  ", "  ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodPost, http.MethodPut:
			var policyReport *v1alpha1.PolicyReport
			raw, _ := io.ReadAll(r.Body)
			h.logger.Info("Body:", string(raw))
			err := json.Unmarshal(raw, &policyReport)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotAcceptable)
			}
			policyReport, err = h.polHandler.Create(ctx, policyReport, vars["namespace"])

			data, err = json.MarshalIndent(policyReport, "  ", "  ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodDelete:
			err = h.polHandler.DeleteCollection(ctx, vars["namespace"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "method not supported on this endpoint", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func (h *handlerSet) NamespacedHandlerWithName(ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		h.logger.Infof("NamespacedHandlerWithName: recieved request method:%s, url:%s, resource:%s, name:%s, namespace:%s", r.Method, r.URL.String(), vars["resource"], vars["name"], vars["namespace"])

		if vars["resource"] != "policyreports" {
			http.Error(w, "only policyreports are supported on this endpoint", http.StatusBadRequest)
		}

		var data []byte = make([]byte, 0)
		var err error
		switch r.Method {
		case http.MethodGet:
			var val *v1alpha1.PolicyReport
			val, err = h.polHandler.Get(ctx, vars["name"], vars["namespace"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotAcceptable)
			}

			data, err = json.MarshalIndent(val, "  ", "  ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodPost, http.MethodPut:
			var policyReport *v1alpha1.PolicyReport
			raw, _ := io.ReadAll(r.Body)
			h.logger.Info("Body:", string(raw))
			err := json.Unmarshal(raw, &policyReport)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotAcceptable)
			}
			policyReport, err = h.polHandler.Update(ctx, policyReport, vars["name"], vars["namespace"])

			data, err = json.MarshalIndent(policyReport, "  ", "  ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodDelete:
			err = h.polHandler.Delete(ctx, vars["name"], vars["namespace"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "method not supported on this endpoint", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
