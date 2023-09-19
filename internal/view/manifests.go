package view

import (
	"context"

	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/tcell/v2"
)

// Helm represents a helm chart view.
type Manifest struct {
	ResourceViewer
}

// NewHelm returns a new alias view.
func NewManifest(gvr client.GVR) ResourceViewer {
	c := Application{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	c.SetContextFn(c.applicationContext)

	return &c
}

func (c *Manifest) applicationContext(ctx context.Context) context.Context {
	return ctx
}
