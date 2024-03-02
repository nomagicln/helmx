// Package postrender provides a way to modify the rendered manifests before they are sent to the Kubernetes cluster.
package postrender

import (
	"bytes"
)

type PostRendererFunc func(renderedManifests *bytes.Buffer) (modifiedManifests *bytes.Buffer, err error)

func (f PostRendererFunc) Run(renderedManifests *bytes.Buffer) (modifiedManifests *bytes.Buffer, err error) {
	return f(renderedManifests)
}
