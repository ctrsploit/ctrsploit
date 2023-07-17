package prerequisite

import (
	"testing"
)

func TestKernelVersion_check(t *testing.T) {
	type fields struct {
		ExpectedMinVersion string
		ExpectedMaxVersion string
		BasePrerequisite   BasePrerequisite
	}
	tests := []struct {
		name          string
		fields        fields
		version       string
		wantSatisfied bool
	}{
		{
			name: "4.6 <= 4.6.0 <= 4.7",
			fields: fields{
				ExpectedMinVersion: "4.6",
				ExpectedMaxVersion: "4.7",
			},
			version:       "4.6.0",
			wantSatisfied: true,
		},
		{
			name: "4.6.0 <= 4.6 <= 4.7",
			fields: fields{
				ExpectedMinVersion: "4.6.0",
				ExpectedMaxVersion: "4.7",
			},
			version:       "4.6",
			wantSatisfied: true,
		},
		{
			name: "4.6 <= 5.15.49 <= ''",
			fields: fields{
				ExpectedMinVersion: "4.6",
				ExpectedMaxVersion: "",
			},
			version:       "5.15.49",
			wantSatisfied: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &KernelVersion{
				ExpectedMinVersion: tt.fields.ExpectedMinVersion,
				ExpectedMaxVersion: tt.fields.ExpectedMaxVersion,
				BasePrerequisite:   tt.fields.BasePrerequisite,
			}
			if gotSatisfied := p.check(tt.version); gotSatisfied != tt.wantSatisfied {
				t.Errorf("check() = %v, want %v", gotSatisfied, tt.wantSatisfied)
			}
		})
	}
}
