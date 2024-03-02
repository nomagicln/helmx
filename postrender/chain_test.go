package postrender

import (
	"bytes"
	"testing"

	"helm.sh/helm/v3/pkg/postrender"
)

func TestChain(t *testing.T) {
	// TestChain tests the Chain post renderer.
	tests := []struct {
		name    string
		input   *bytes.Buffer
		prs     []postrender.PostRenderer
		want    string
		wantErr bool
	}{
		{
			name:  "empty",
			input: bytes.NewBufferString("Seed"),
			prs:   nil,
			want:  "Seed",
		},
		{
			name:  "single",
			input: bytes.NewBufferString("Seed"),
			prs: []postrender.PostRenderer{
				PostRendererFunc(func(renderedManifests *bytes.Buffer) (*bytes.Buffer, error) {
					return bytes.NewBuffer(append(renderedManifests.Bytes(), []byte(":Foo")...)), nil
				}),
			},
			want: "Seed:Foo",
		},
		{
			name:  "multiple",
			input: bytes.NewBufferString("Seed"),
			prs: []postrender.PostRenderer{
				PostRendererFunc(func(renderedManifests *bytes.Buffer) (*bytes.Buffer, error) {
					return bytes.NewBuffer(append(renderedManifests.Bytes(), []byte(":Foo")...)), nil
				}),
				PostRendererFunc(func(renderedManifests *bytes.Buffer) (*bytes.Buffer, error) {
					return bytes.NewBuffer(append(renderedManifests.Bytes(), []byte(":Bar")...)), nil
				}),
			},
			want: "Seed:Foo:Bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := NewChain(tt.prs...)
			got, err := chain.Run(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chain.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("Chain.Run() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}
