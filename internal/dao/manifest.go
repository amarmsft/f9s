package dao

import (
	"context"
	"fmt"

	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/render"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	_ Accessor = (*Application)(nil)
	_ Nuker    = (*Application)(nil)
)

// CustomResourceDefinition represents a CRD resource model.
type Manifest struct {
	NonResource
}

// List returns a collection of nodes.
func (c *Manifest) List(ctx context.Context, _ string) ([]runtime.Object, error) {
	fqn, ok := ctx.Value(internal.KeyPath).(string)
	if !ok {
		return nil, fmt.Errorf("no context path for %q", c.gvr)
	}

	//fqn = "clusterfleet/" + fqn
	app, err := c.fetchApplication(fqn)
	if err != nil {
		return nil, err
	}
	res := make([]runtime.Object, 0, len(app.Spec.Workload))
	statuses := make(map[string]render.ManifestStatus, 0)
	for _, status := range app.Status.ManifestStatuses {
		statuses[fmt.Sprintf("%s.%s.%s", status.Namespace, status.Name, status.Kind)] = status
	}

	for _, co := range app.Spec.Workload {
		res = append(res, c.makeManifestResp(co.ManifestItem, statuses))
	}

	return res, nil
}

func (c *Manifest) fetchApplication(fqn string) (*render.Application, error) {
	o, err := c.GetFactory().Get("apis.clusterfleet.io/v1alpha1/applications", fqn, true, labels.Everything())
	if err != nil {
		return nil, err
	}
	var app render.Application
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(o.(*unstructured.Unstructured).Object, &app)
	return &app, err
}

func (c *Manifest) makeManifestResp(manifest render.Manifest, statuses map[string]render.ManifestStatus) render.ManifestRes {

	mo, _ := render.ManifestToUnstructed(manifest)

	name := render.GetNameFromUnstructured(mo)
	namespace := render.GetNamespaceFromUnstructured(mo)
	kind := render.GetKindFromUnstructured(mo)
	replicas := render.GetReplicaForUnstructured(mo)
	status, _ := statuses[fmt.Sprintf("%s.%s.%s", namespace, name, kind)]

	return render.ManifestRes{
		Manifest:  &manifest,
		Status:    &status,
		Name:      name + "_" + kind,
		Namespace: namespace,
		Kind:      kind,
		Replicas:  replicas,
	}
}
