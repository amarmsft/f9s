package dao

import (
	"context"
	"fmt"

	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/render"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	_ Accessor = (*Application)(nil)
	_ Nuker    = (*Application)(nil)
)

// CustomResourceDefinition represents a CRD resource model.
type ApplicationStatus struct {
	NonResource
}

// List returns a collection of nodes.
func (c *ApplicationStatus) List(ctx context.Context, _ string) ([]runtime.Object, error) {
	fqn, ok := ctx.Value(internal.KeyPath).(string)
	if !ok {
		return nil, fmt.Errorf("no context path for %q", c.gvr)
	}

	//fqn = "clusterfleet/" + fqn
	app, err := c.fetchApplication(fqn)
	if err != nil {
		return nil, err
	}
	res := make([]runtime.Object, 0, len(app.Status.Clusters))

	for _, w := range app.Status.Clusters {
		res = append(res, c.makeApplicationStatusResp(w))
	}

	return res, nil
}

func (c *ApplicationStatus) fetchApplication(fqn string) (*render.Application, error) {
	o, err := c.GetFactory().Get("apis.clusterfleet.io/v1alpha1/applications", fqn, true, labels.Everything())
	if err != nil {
		return nil, err
	}
	var app render.Application
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(o.(*unstructured.Unstructured).Object, &app)
	return &app, err
}

func (c *ApplicationStatus) makeApplicationStatusResp(cluster render.ApplicationClusterStatus) render.ApplicationStatusRes {
	var readyCondition v1.Condition = v1.Condition{
		Reason: "Unknown",
	}

	for _, cond := range cluster.Conditions {
		if cond.Type == "Ready" {
			readyCondition = cond
		}
	}

	ready := readyCondition.Reason

	return render.ApplicationStatusRes{
		Cluster: cluster.Cluster,
		Ready:   ready,
		Status:  &cluster,
	}
}
