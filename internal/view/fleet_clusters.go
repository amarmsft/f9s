package view

import (
	"context"

	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/tcell/v2"
)

// Helm represents a helm chart view.
type FleetClusters struct {
	ResourceViewer
}

// NewHelm returns a new alias view.
func NewFleetCluster(gvr client.GVR) ResourceViewer {
	c := FleetClusters{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.SetContextFn(c.clustersContext)

	return &c
}

func (c *FleetClusters) clustersContext(ctx context.Context) context.Context {
	key := c.GetTable().GetSelectedItem()
	return context.WithValue(ctx, internal.KeyPath, key)
}
