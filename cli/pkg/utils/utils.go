package utils

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func IsPolicyReport(resource string) bool {
	return resource == "polr" || resource == "policyreport" || resource == "policyreports"
}

func IsClusterPolicyReport(resource string) bool {
	return resource == "cpolr" || resource == "clusterpolicyreport" || resource == "clusterpolicyreports"
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
