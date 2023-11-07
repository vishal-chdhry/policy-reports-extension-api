package utils

import (
	"encoding/json"

	"github.com/vishal-chdhry/policy-reports-extension-api/client/pkg/v1alpha1"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func IsPolicyReport(resource string) bool {
	return resource == "polr" || resource == "policyreport" || resource == "policyreports"
}

func IsClusterPolicyReport(resource string) bool {
	return resource == "cpolr" || resource == "clusterpolicyreport" || resource == "clusterpolicyreports"
}

func PolicyReportToUnstructured(pol *v1alpha1.PolicyReport) *unstructured.Unstructured {
	b, err := json.Marshal(pol)
	if err != nil {
		panic(err)
	}
	intr := bytesToInterface(b)

	return &unstructured.Unstructured{
		Object: intr,
	}
}

func ClusterPolicyReportToUnstructured(cpol *v1alpha1.ClusterPolicyReport) *unstructured.Unstructured {
	b, err := json.Marshal(cpol)
	if err != nil {
		panic(err)
	}
	intr := bytesToInterface(b)

	return &unstructured.Unstructured{
		Object: intr,
	}
}

func ClusterPolicyReportListToUnstructuredList(cpol *v1alpha1.ClusterPolicyReportList) *unstructured.UnstructuredList {
	b, err := json.Marshal(cpol)
	if err != nil {
		panic(err)
	}
	intr := bytesToInterface(b)

	ul := make([]unstructured.Unstructured, 0)
	for _, v := range cpol.Items {
		unst := ClusterPolicyReportToUnstructured(&v)
		ul = append(ul, *unst)
	}

	return &unstructured.UnstructuredList{
		Object: intr,
		Items:  ul,
	}
}

func PolicyReportListToUnstructuredList(pol *v1alpha1.PolicyReportList) *unstructured.UnstructuredList {
	b, err := json.Marshal(pol)
	if err != nil {
		panic(err)
	}
	intr := bytesToInterface(b)

	ul := make([]unstructured.Unstructured, 0)
	for _, v := range pol.Items {
		unst := PolicyReportToUnstructured(&v)
		ul = append(ul, *unst)
	}

	return &unstructured.UnstructuredList{
		Object: intr,
		Items:  ul,
	}
}

func bytesToInterface(b []byte) map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal(b, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func JSONToUnstructured(b []byte) *unstructured.Unstructured {
	obj := make(map[string]interface{})
	err := json.Unmarshal(b, &obj)
	if err != nil {
		panic(err)
	}

	return &unstructured.Unstructured{
		Object: obj,
	}
}

func YAMLToUnstructured(b []byte) *unstructured.Unstructured {
	obj := make(map[string]interface{})
	err := yaml.Unmarshal(b, &obj)
	if err != nil {
		panic(err)
	}

	return &unstructured.Unstructured{
		Object: obj,
	}
}
