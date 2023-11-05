# policy-reports-extension-api

Policy Reports Extension
===============================================

This project is a Kubernetes [API
Extension](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/apiserver-aggregation/)
to store Policy reports and cluster policy reports in an external database instead of etcd.

The etcd database currently has a maximum size of 8GB. Reports tend to be relatively large objects, but even with very small reports, the etcd capacity can be easily reached in larger clusters with many report producers. Under heavy report activity (e.g. cluster churn, scanning, analytical processes, etc.), the volume of data being written and retrieved by etcd requires the API server to buffer large amounts of data. This compounds existing cluster issues and can cascade into complete API unavailability. CAP guarantees are not required for reports, which, at present, are understood to be ephemeral data which will be re-created if deleted. Philosophically, report data is analytical in nature and should not be stored in the transactional database. These are also the arguments for [removing Kubernetes `Events` from the primary etcd](https://github.com/kubernetes/kubernetes/issues/4432) as well. This issue has not yet been fully resolved.

This extension allows us store policy reports outside of etcd.

Build
-----

Build the docker image for the server using `server/Dockerfile`

Apply the manifests in a kubernetes cluster:

```
kubectl apply -f config/manifest.yaml
```

Usage
-----

// TODO
