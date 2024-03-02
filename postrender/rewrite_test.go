package postrender

import (
	"bytes"
	"testing"
)

func TestRewriter(t *testing.T) {
	manifest := `apiVersion: v1
kind: Pod
metadata:
  name: test
---
apiVersion: v1
kind: Deployment
metadata:
  name: test
`
	tests := []struct {
		name     string
		options  []RewriteOption
		manifest string
		expected string
	}{
		{
			name:     "enhance without option",
			manifest: manifest,
			expected: manifest,
		},
		{
			name: "enhance with option",
			options: []RewriteOption{
				AppendLabels(map[string]string{"app": "test"}),
				AppendAnnotations(map[string]string{"app": "test"}),
			},
			manifest: manifest,
			expected: `apiVersion: v1
kind: Pod
metadata:
  annotations:
    app: test
  labels:
    app: test
  name: test
---
apiVersion: v1
kind: Deployment
metadata:
  annotations:
    app: test
  labels:
    app: test
  name: test
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifiedManifests, err := NewRewriter(tt.options).Run(bytes.NewBufferString(tt.manifest))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got := modifiedManifests.String(); got != tt.expected {
				t.Errorf("expected: %s, got: %s", tt.expected, got)
			}
		})
	}
}
