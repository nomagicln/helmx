package postrender

import (
	"bytes"
	"errors"
	"io"

	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/postrender"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// RewriteOption is a function type that represents an option for enhancing an unstructured object.
type RewriteOption func(unstructured.Unstructured)

type rewriter struct {
	options []RewriteOption
}

func (e *rewriter) Run(renderedManifests *bytes.Buffer) (*bytes.Buffer, error) {
	var (
		modifiedManifests = new(bytes.Buffer)
		decoder           = yaml.NewDecoder(renderedManifests)
		encoder           = yaml.NewEncoder(modifiedManifests)
	)

	encoder.SetIndent(2) // Set the indent to 2 spaces, as the default is 4 spaces.

	for {
		var object unstructured.Unstructured
		if err := decoder.Decode(&object.Object); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		for _, option := range e.options {
			option(object)
		}

		// Unfortunately, the yaml style will be changed after encoding. And the comments will be lost.
		if err := encoder.Encode(object.Object); err != nil {
			return nil, err
		}
	}

	return modifiedManifests, nil
}

// NewRewriter applies a list of EnhanceOption functions to modify the renderedManifests.
// It returns a PostRenderer that applies the enhancements to the manifests.
// The EnhanceOption functions are applied in the order they appear in the options slice.
// Each EnhanceOption function takes an unstructured.Unstructured object as input and modifies it.
// The modified manifests are encoded and returned as a *bytes.Buffer.
// If any error occurs during the enhancement process, it is returned.
func NewRewriter(options []RewriteOption) postrender.PostRenderer {
	return &rewriter{options: options}
}

func mergeMap[K comparable, V any](dicts ...map[K]V) map[K]V {
	var m = make(map[K]V)

	for _, dict := range dicts {
		for k, v := range dict {
			m[k] = v
		}
	}

	return m
}

func AppendLabels(labels map[string]string) RewriteOption {
	return func(obj unstructured.Unstructured) {
		obj.SetLabels(mergeMap(obj.GetLabels(), labels))
	}
}

func AppendAnnotations(annotations map[string]string) RewriteOption {
	return func(obj unstructured.Unstructured) {
		obj.SetAnnotations(mergeMap(obj.GetAnnotations(), annotations))
	}
}
