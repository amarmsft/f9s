package dao

import (
	"context"

	"github.com/derailed/k9s/internal"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	_ Accessor = (*Application)(nil)
	_ Nuker    = (*Application)(nil)
)

// CustomResourceDefinition represents a CRD resource model.
type Application struct {
	Resource
}

// List returns a collection of nodes.
func (c *Application) List(ctx context.Context, _ string) ([]runtime.Object, error) {
	strLabel, ok := ctx.Value(internal.KeyLabels).(string)
	labelSel := labels.Everything()
	if sel, e := labels.ConvertSelectorToLabelsMap(strLabel); ok && e == nil {
		labelSel = sel.AsSelector()
	}

	const gvr = "apis.clusterfleet.io/v1alpha1/applications"
	return c.GetFactory().List(gvr, "-", false, labelSel)
}
