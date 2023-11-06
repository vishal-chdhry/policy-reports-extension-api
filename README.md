Policy Reports Extension (PRExt)
===============================================

*Inspired by https://github.com/cmurphy/hns-list*

This project is a Kubernetes [API
Extension](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/apiserver-aggregation/)
to store Policy reports and cluster policy reports in an external database instead of etcd.

Why move policy reports outside of etcd?: https://github.com/kyverno/KDP/pull/51

> It is desirable to move reports out of etcd for several reasons:
> 
> Why leave etcd:
> 
> - The etcd database currently has a maximum size of 8GB. Reports tend to be relatively large objects, but even with very small reports, the etcd capacity can be easily reached in larger clusters with many report producers.
> - Under heavy report activity (e.g. cluster churn, scanning, analytical processes, etc.), the volume of data being written and retrieved by etcd requires the API server to buffer large amounts of data. This compounds existing cluster issues and can cascade into complete API unavailability.
> - CAP guarantees are not required for reports, which, at present, are understood to be ephemeral data which will be re-created if deleted.
> - Philosophically, report data is analytical in nature and should not be stored in the transactional database.
> - These are also the arguments for [removing Kubernetes `Events` from the primary etcd](https://github.com/kubernetes/kubernetes/issues/4432) as well. This issue has not yet been fully resolved.
> 
> Expected benefits of alternative storage:
> 
> - Alleviation of the etcd + API server load and capacity limitations.
> - Common report consumer workflows can be more efficient.
>     - Report consumers are often analytical in nature and/or operate on aggregate data. The API is not designed to efficiently handle, for example, a query for all reports where containing a vulnerability with a CVSS severity of 8.0 or above. To perform such a query, a report consumer must retrieve and parse all of the reports. Retrieving a large volume of reports, especially with multiple simultaneous consumers, leads to the performance issues described previously.
>     - With reports stored in, for instance, a relational database, report consumers could instead query the underlying database directly, using more robust query syntax.
>     - This would improve the implementation of, or even replace the need for, certain exporters, and enable new reporting use cases.


This extension allows us store policy reports outside of etcd.

Build
-----

Build the kubectl-prext plugin and install it in your GOPATH:

```
make cli
```

Build the docker image for the server:

```
make server
```

Apply the manifests in a kubernetes cluster:

```
kubectl apply -f config/manifest.yaml
```

Usage
-----

// TODO
