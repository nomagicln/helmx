package postrender

import (
	"bytes"

	"helm.sh/helm/v3/pkg/postrender"
)

// chain is a post renderer that runs a list of post renderers in order.
type chain struct {
	// The list of post renderers to run in order.
	prs []postrender.PostRenderer
}

// NewChain creates a new Chain post renderer that runs the given list of post renderers in order.
func NewChain(prs ...postrender.PostRenderer) postrender.PostRenderer {
	return &chain{prs: prs}
}

// Run implements postrender.PostRenderer.
func (c *chain) Run(renderedManifests *bytes.Buffer) (*bytes.Buffer, error) {
	var (
		modifiedManifests *bytes.Buffer = renderedManifests
		err               error
	)

	for _, pr := range c.prs {
		if modifiedManifests, err = pr.Run(modifiedManifests); err != nil {
			return nil, err
		}
	}
	return modifiedManifests, nil
}
