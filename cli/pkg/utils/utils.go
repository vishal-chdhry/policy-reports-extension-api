package utils

import (
	"reflect"

	"github.com/vishal-chdhry/policy-reports-extension-api/client/pkg/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func IsPolicyReport(resource string) bool {
	return resource == "polr" || resource == "policyreport" || resource == "policyreports"
}

func IsClusterPolicyReport(resource string) bool {
	return resource == "cpolr" || resource == "clusterpolicyreport" || resource == "clusterpolicyreports"
}

func PolicyReportToUnstructured(pol *v1alpha1.PolicyReport) *unstructured.Unstructured {
	intr := reflect.ValueOf(pol).Interface().(map[string]interface{})

	return &unstructured.Unstructured{
		Object: intr,
	}
}

func ClusterPolicyReportToUnstructured(cpol *v1alpha1.ClusterPolicyReport) *unstructured.Unstructured {
	intr := reflect.ValueOf(cpol).Interface().(map[string]interface{})

	return &unstructured.Unstructured{
		Object: intr,
	}
}

func ClusterPolicyReportListToUnstructuredList(cpol *v1alpha1.ClusterPolicyReportList) *unstructured.UnstructuredList {
	intr := reflect.ValueOf(cpol).Interface().(map[string]interface{})

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

func PolicyReportListToUnstructuredList(cpol *v1alpha1.PolicyReportList) *unstructured.UnstructuredList {
	intr := reflect.ValueOf(cpol).Interface().(map[string]interface{})

	ul := make([]unstructured.Unstructured, 0)
	for _, v := range cpol.Items {
		unst := PolicyReportToUnstructured(&v)
		ul = append(ul, *unst)
	}

	return &unstructured.UnstructuredList{
		Object: intr,
		Items:  ul,
	}
}
