package main

import (
	"reflect"
	"testing"
)

func TestMaintainersFromReadme(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
		wantErr bool
	}{
		{
			name: "scalar maintainer",
			content: `---
maintainer: ssst0n3
---
# Module
`,
			want: []string{"ssst0n3"},
		},
		{
			name: "list maintainer",
			content: `---
maintainer:
  - r0binak
  - ssst0n3
---
# Module
`,
			want: []string{"r0binak", "ssst0n3"},
		},
		{
			name: "trims empty maintainers",
			content: `---
maintainer:
  - ""
  - " ssst0n3 "
---
# Module
`,
			want: []string{"ssst0n3"},
		},
		{
			name: "ignores author",
			content: `---
author: ssst0n3
---
# Module
`,
			want: nil,
		},
		{
			name: "no front matter",
			content: `# Module
`,
			want: nil,
		},
		{
			name: "invalid yaml",
			content: `---
maintainer: [
---
# Module
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := maintainersFromReadme([]byte(tt.content))
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("maintainers = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestResolveMaintainersFallsBackToDefault(t *testing.T) {
	got := resolveMaintainers(t.TempDir(), "")
	want := []string{defaultMaintainer}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("maintainers = %#v, want %#v", got, want)
	}
}
