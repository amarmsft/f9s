package render

import (
	"encoding/json"
	"fmt"

	"github.com/derailed/k9s/internal/client"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// CustomResourceDefinition renders a K8s CustomResourceDefinition to screen.
type ApplicationStatusRenderer struct {
	Base
}

// Header returns a header rbw.
func (ApplicationStatusRenderer) Header(string) Header {
	return Header{
		HeaderColumn{Name: "CLUSTER"},
		HeaderColumn{Name: "READY"},
		HeaderColumn{Name: "REPLICAS"},
	}
}

// Render renders a K8s resource to screen.
func (c ApplicationStatusRenderer) Render(o interface{}, ns string, r *Row) error {
	appStatus, ok := o.(ApplicationStatusRes)
	if !ok {
		return fmt.Errorf("Expected ManifestRes, but got %T", o)
	}

	replicas := "-"
	for _, status := range appStatus.Status.ManifestStatuses {
		if status.Kind == "Deployment" {
			deploymentStatus, _ := c.DeploymentStatusToMap(status)
			repl, _, _ := unstructured.NestedFloat64(*deploymentStatus, "replicas")
			avl, found, _ := unstructured.NestedFloat64(*deploymentStatus, "availableReplicas")

			if found {
				replicas = fmt.Sprint(avl) + "/" + fmt.Sprint(repl)
			} else {
				replicas = "0/" + fmt.Sprint(repl)
			}
		}
	}

	r.ID = client.FQN(client.ClusterScope, appStatus.Cluster)
	r.Fields = Fields{
		appStatus.Cluster,
		appStatus.Ready,
		replicas,
	}

	return nil
}

func (c ApplicationStatusRenderer) DeploymentStatusToMap(
	manifestStatus ManifestStatus,
) (*map[string]interface{}, error) {
	var status map[string]interface{}
	err := json.Unmarshal(manifestStatus.Status.Raw, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

type ApplicationStatusRes struct {
	Cluster string
	Ready   string
	Status  *ApplicationClusterStatus
	Error   string
}

// GetObjectKind returns a schema object.
func (c ApplicationStatusRes) GetObjectKind() schema.ObjectKind {
	return nil
}

// DeepCopyObject returns a container copy.
func (c ApplicationStatusRes) DeepCopyObject() runtime.Object {
	return c
}
