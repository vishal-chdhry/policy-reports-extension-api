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

Build the docker image for the server:

```sh
make server
```

Install
-------

1. Install [cert-manager](https://cert-manager.io/docs/installation/)

2. Apply the manifest:

```sh
kubectl apply -f manifest/manifest.yaml
```

Usage
-----

Create a new policy report:

```sh
kubectl create --raw /apis/prext.demo/v1alpha1/namespaces/default/policyreports -f config/testdata/testpolicy.json
```

Output
```sh
{
    "kind": "PolicyReport",
    "apiVersion": "prext.demo/v1alpha1",
    "metadata": {
      "name": "test",
      "namespace": "default",
      "uid": "8f611e66-e398-401d-9624-fd9ec87a72f8",
      "creationTimestamp": null,
      "ownerReferences": [
        {
          "apiVersion": "v1",
          "kind": "Pod",
...
...
```

Get the all policy reports in a namespace:

```sh
kubectl get --raw /apis/prext.demo/v1alpha1/namespaces/{{NAMESPACE}}/policyreports | jq --args ".[].metadata.name"      
```

Output: 
```sh
{
    "kind": "PolicyReportList",
    "apiVersion": "prext.demo/v1alpha1",
    "metadata": {
      "resourceVersion": "31622"
    },
    "items": [
      {
        "kind": "PolicyReport",
        "apiVersion": "prext.demo/v1alpha1",
        "metadata": {
          "name": "image-scan-pod-identity-webhook-6d5f85488d-2vdjn",
...
...
```

View a policy report:
```sh
kubectl get --raw /apis/prext.demo/v1alpha1/namespaces/default/policyreports/test1
```

Output:
```sh
{
    "kind": "PolicyReport",
    "apiVersion": "prext.demo/v1alpha1",
    "metadata": {
      "name": "image-scan-pod-identity-webhook-6d5f85488d-2vdjn",
      "namespace": "default",
      "uid": "87428580-1dbf-4f2c-8320-069c9db00367",
      "resourceVersion": "125248938",
      "generation": 19,
      "creationTimestamp": "2023-10-18T15:33:08Z",
      "labels": {
        "kubernetes.io/createdBy": "vulnerability-operator"
...
...
```

Delete a policy report: 
```sh
kubectl delete --raw /apis/prext.demo/v1alpha1/namespaces/default/policyreports/test 
```

View a cluster policy report: 
```sh
kubectl get --raw /apis/prext.demo/v1alpha1/clusterpolicyreports/test1 
```

Output:
```sh
{
    "kind": "ClusterPolicyReport",
    "apiVersion": "prext.demo/v1alpha1",
    "metadata": {
      "name": "kube-bench",
      "uid": "ed0d888c-fab5-4597-98f4-d93bbd2a28d6",
      "resourceVersion": "130199899",
      "generation": 3,
      "creationTimestamp": "2023-10-22T00:00:14Z",
      "managedFields": [
        {
          "manager": "policyreport",
          "operation": "Update",
          "apiVersion": "prext.demo/v1alpha1",
          "time": "2023-10-22T00:00:14Z",
          "fieldsType": "FieldsV1",
          "fieldsV1": {
            "f:results": {},
            "f:summary": {
              ".": {},
              "f:error": {},
              "f:fail": {},
              "f:pass": {},
              "f:skip": {},
              "f:warn": {}
            }
          }
        }
      ]
    },
    "summary": {
      "pass": 14,
      "fail": 0,
      "warn": 2,
...
...
```

Get all cluster policy reports: 
```sh
kubectl get --raw /apis/prext.demo/v1alpha1/clusterpolicyreports
```

Output:
```sh
{
    "kind": "ClusterPolicyReportList",
    "apiVersion": "prext.demo/v1alpha1",
    "metadata": {
      "resourceVersion": "31778"
    },
    "items": [
      {
        "kind": "ClusterPolicyReport",
        "apiVersion": "prext.demo/v1alpha1",
        "metadata": {
          "name": "kube-bench",
          "uid": "ed0d888c-fab5-4597-98f4-d93bbd2a28d6",
          "resourceVersion": "130199899",
          "generation": 3,
          "creationTimestamp": "2023-10-22T00:00:14Z",
          "managedFields": [
            {
              "manager": "policyreport",
              "operation": "Update",
              "apiVersion": "prext.demo/v1alpha1",
...
...
```