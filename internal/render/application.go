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
type ApplicationRenderer struct {
	Base
}

// ColorerFunc colors a resource row.
func (ApplicationRenderer) ColorerFunc() ColorerFunc {
	return func(ns string, h Header, re RowEvent) tcell.Color {
		c := DefaultColorer(ns, h, re)

		statusCol := h.IndexOf("PROVISIONED", true)
		if statusCol == -1 {
			return c
		}
		status := strings.TrimSpace(re.Row.Fields[statusCol])
		switch status {
		case "ApplicationProvisioned":
			c = StdColor
		default:
			c = ErrColor
		}

		return c
	}
}

// Header returns a header rbw.
func (ApplicationRenderer) Header(string) Header {
	return Header{
		HeaderColumn{Name: "NAMESPACE"},
		HeaderColumn{Name: "NAME"},
		HeaderColumn{Name: "PROVISIONED"},
		HeaderColumn{Name: "CLUSTERS"},
		HeaderColumn{Name: "AGE", Time: true},
	}
}

// Render renders a K8s resource to screen.
func (c ApplicationRenderer) Render(o interface{}, ns string, r *Row) error {
	raw, ok := o.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("Expected CustomResourceDefinition, but got %T", o)
	}

	var app Application
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw.Object, &app)
	if err != nil {
		return err
	}

	provisioned := "Not Processed"
	for _, cond := range app.Status.Conditions {
		if cond.Type == "Provisioned" {
			provisioned = cond.Reason
		}
	}

	clustersToShow := make([]string, 0, 30)
	for _, clusters := range app.Status.Clusters {
		clustersToShow = append(clustersToShow, clusters.Cluster)
	}

	r.ID = client.FQN(app.Namespace, app.GetName())
	r.Fields = Fields{
		app.GetNamespace(),
		app.GetName(),
		provisioned,
		Truncate(join(clustersToShow, ","), 30),
		toAge(app.GetCreationTimestamp()),
	}

	return nil
}
