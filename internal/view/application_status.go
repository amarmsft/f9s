package view

import (
	"context"

	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/tcell/v2"
)

// Helm represents a helm chart view.
type ApplicationStatus struct {
	ResourceViewer
}

// NewHelm returns a new alias view.
func NewApplicationStatus(gvr client.GVR) ResourceViewer {
	c := ApplicationStatus{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.SetContextFn(c.applicationContext)

	return &c
}

func (c *ApplicationStatus) applicationContext(ctx context.Context) context.Context {
	key := c.GetTable().GetSelectedCell(0) + "/" + c.GetTable().GetSelectedCell(1)
	return context.WithValue(ctx, internal.KeyPath, key)
}

// Name returns the component name.
func (c *ApplicationStatus) Name() string { return "status" }
