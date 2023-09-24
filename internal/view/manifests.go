package view

import (
	"context"
	"time"

	"github.com/derailed/k9s/internal/client"
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

	c.GetTable().GetModel().SetRefreshRate(10 * time.Minute)
	//c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	//c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Attributes(tcell.AttrNone))
	c.SetContextFn(c.applicationContext)

	return &c
}

func (c *Manifest) applicationContext(ctx context.Context) context.Context {
	return ctx
}

// Name returns the component name.
func (c *Manifest) Name() string { return "manifests" }
