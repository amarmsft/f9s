package render

import (
	"fmt"
	"strings"

	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/tcell/v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// CustomResourceDefinition renders a K8s CustomResourceDefinition to screen.
type ClusterRenderer struct {
	Base
}

// ColorerFunc colors a resource row.
func (ClusterRenderer) ColorerFunc() ColorerFunc {
	return func(ns string, h Header, re RowEvent) tcell.Color {
		c := DefaultColorer(ns, h, re)

		statusCol := h.IndexOf("PROVISIONED", true)
		if statusCol == -1 {
			return c
		}
		status := strings.TrimSpace(re.Row.Fields[statusCol])
		switch status {
		case "Ready":
			c = StdColor
		case "Provisioning":
			c = PendingColor
		default:
			c = ErrColor
		}

		return c
	}
}

// Header returns a header rbw.
func (ClusterRenderer) Header(string) Header {
	return Header{
		HeaderColumn{Name: "NAME"},
		HeaderColumn{Name: "USGAE"},
		HeaderColumn{Name: "PROVISIONED"},
		HeaderColumn{Name: "HEALTH"},
		HeaderColumn{Name: "AGE", Time: true},
	}
}

// Render renders a K8s resource to screen.
func (c ClusterRenderer) Render(o interface{}, ns string, r *Row) error {
	raw, ok := o.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("Expected CustomResourceDefinition, but got %T", o)
	}

	var cl Cluster
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw.Object, &cl)
	if err != nil {
		return err
	}

	provisioned := "Not Processed"
	for _, st := range cl.Status.ProvisioningStatus {
		if st.Name == "provisioner" {
			provisioned = st.Status
		}
	}

	r.ID = client.FQN(client.ClusterScope, cl.GetName())
	r.Fields = Fields{
		cl.GetName(),
		"general",
		provisioned,
		string(cl.Status.ClusterHealthStatus),
		toAge(cl.GetCreationTimestamp()),
	}

	return nil
}
