package view

import (
	"context"

	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/ui"
	"github.com/derailed/tcell/v2"
)

// Helm represents a helm chart view.
type Application struct {
	ResourceViewer
}

// NewHelm returns a new alias view.
func NewApplication(gvr client.GVR) ResourceViewer {
	c := Application{
		ResourceViewer: NewBrowser(gvr),
	}
	c.GetTable().SetBorderFocusColor(tcell.ColorMediumSpringGreen)
	c.GetTable().SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorMediumSpringGreen).Attributes(tcell.AttrNone))
	//c.SetContextFn(c.applicationContext)
	c.GetTable().SetEnterFn(c.showManifests)
	c.AddBindKeysFn(c.bindKeys)

	return &c
}

func (c *Application) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyS: ui.NewKeyAction("Show Status", c.showApplicationStatus, true),
	})
	aa.Add(resourceSorters(c.GetTable()))
}

func (c *Application) showManifests(app *App, model ui.Tabular, gvr, path string) {
	co := NewManifest(client.NewGVR("manifests"))
	co.SetContextFn(c.applicationContext(path))
	if err := app.inject(co, false); err != nil {
		app.Flash().Err(err)
	}
}

func (c *Application) showApplicationStatus(evt *tcell.EventKey) *tcell.EventKey {
	path := c.GetTable().GetSelectedItem()
	if path == "" {
		return evt
	}

	status := NewApplicationStatus(client.NewGVR("applicationStatus"))
	status.SetContextFn(c.applicationContext(path))
	//no.SetContextFn(nodeContext(pod.Spec.NodeName))
	if err := c.App().inject(status, false); err != nil {
		c.App().Flash().Err(err)
	}

	return nil
}

func (c *Application) applicationContext(path string) ContextFunc {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, internal.KeyPath, path)
	}
}

/*func (c *Application) applicationContext(ctx context.Context) context.Context {
	key := c.GetTable().GetSelectedCell(0) + "/" + c.GetTable().GetSelectedCell(1)
	return context.WithValue(ctx, internal.KeyPath, key)
}*/
