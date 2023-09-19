package render

import (
	"encoding/json"
	"fmt"

	"github.com/derailed/k9s/internal/client"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

// CustomResourceDefinition renders a K8s CustomResourceDefinition to screen.
type ManifestRenderer struct {
	Base
}

// Header returns a header rbw.
func (ManifestRenderer) Header(string) Header {
	return Header{
		HeaderColumn{Name: "NAME"},
		HeaderColumn{Name: "KIND"},
		HeaderColumn{Name: "REPLICAS"},
	}
}

// Render renders a K8s resource to screen.
func (c ManifestRenderer) Render(o interface{}, ns string, r *Row) error {
	manifest, ok := o.(ManifestRes)
	if !ok {
		return fmt.Errorf("Expected ManifestRes, but got %T", o)
	}

	replicas := "-"
	if manifest.Kind == "Deployment" {
		deploymentStatus, _ := c.DeploymentStatusToMap(*manifest.Status)
		repl, _, _ := unstructured.NestedFloat64(*deploymentStatus, "replicas")
		avl, found, _ := unstructured.NestedFloat64(*deploymentStatus, "availableReplicas")

		if found {
			replicas = fmt.Sprint(avl) + "/" + fmt.Sprint(repl)
		} else {
			replicas = "0/" + fmt.Sprint(repl)
		}
	}

	r.ID = client.FQN(client.ClusterScope, manifest.Name)
	r.Fields = Fields{
		manifest.Name,
		manifest.Kind,
		replicas,
	}

	return nil
}

func (c ManifestRenderer) DeploymentStatusToMap(
	manifestStatus ManifestStatus,
) (*map[string]interface{}, error) {
	var status map[string]interface{}
	err := json.Unmarshal(manifestStatus.Status.Raw, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func GetKindFromUnstructured(unstructuredObj *unstructured.Unstructured) string {
	value, found, err := unstructured.NestedString(unstructuredObj.Object, "kind")
	if !found || err != nil {
		return ""
	}
	return string(value)
}

func GetNameFromUnstructured(unstructuredObj *unstructured.Unstructured) string {
	value, found, err := unstructured.NestedString(unstructuredObj.Object, "metadata", "name")
	if !found || err != nil {
		return ""
	}
	return string(value)
}

func GetNamespaceFromUnstructured(unstructuredObj *unstructured.Unstructured) string {
	value, found, err := unstructured.NestedString(unstructuredObj.Object, "metadata", "namespace")
	if !found || err != nil {
		return ""
	}
	return string(value)
}

func GetReplicaForUnstructured(unstructuredObj *unstructured.Unstructured) string {
	value, found, err := unstructured.NestedInt64(unstructuredObj.Object, "spec", "replicas")
	if !found || err != nil {
		return "-"
	}
	return fmt.Sprint(value)
}

func ManifestToUnstructed(manifest Manifest) (*unstructured.Unstructured, error) {
	// If manifest was applied through YAML document, it will be populated in manifest.Raw.
	// But if applied through scheduler, it is possible that manifest.Raw is empty and
	// data is in manifest.Object instead. Webhook needs to be able to convert in both cases.
	if manifest.Raw != nil {
		ret, err := getUnstructuredFromRawExtension(&manifest.RawExtension)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to convert manifest.Raw. raw: %s, error: %w",
				string(manifest.Raw),
				err,
			)
		}
		return ret, nil
	} else {
		obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(manifest.Object)
		if err != nil {
			return nil, fmt.Errorf("failed to convert manifest.Object. object: %v, error: %w", manifest.Object, err)
		}
		return &unstructured.Unstructured{Object: obj}, nil
	}
}

func getUnstructuredFromRawExtension(
	raw *runtime.RawExtension,
) (*unstructured.Unstructured, error) {
	if raw == nil {
		return nil, fmt.Errorf(`unable to create runtime object from raw extension: nil`)
	}

	deserializer := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj, _, err := deserializer.Decode(raw.Raw, nil, nil)
	if err != nil {
		return nil, fmt.Errorf(
			`unable to convert raw extension %v into runtime object: %w`,
			*raw,
			err,
		)
	}

	return getUnstructuredFromObject(obj)
}

func getUnstructuredFromObject(
	obj interface{},
) (*unstructured.Unstructured, error) {
	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, fmt.Errorf(
			`unable to convert runtime object %v into unstructured map: %w`,
			obj,
			err,
		)
	}

	return &unstructured.Unstructured{Object: unstructuredMap}, nil
}

type ManifestRes struct {
	Manifest  *Manifest
	Name      string
	Namespace string
	Kind      string
	Replicas  string
	Status    *ManifestStatus
}

// GetObjectKind returns a schema object.
func (c ManifestRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (c ManifestRes) DeepCopyObject() runtime.Object {
	return c
}
