package status

import (
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		want    *Status
		wantErr bool
	}{
		{
			name: "golden path",
			content: `
Name:   my-process
Uid:    1000    1001    1002    1003
Gid:    2000    2001    2002    2003
CapInh: 00000000a80425fb
CapPrm: 00000000a80425fb
CapEff: 00000000a80425fb
CapBnd: 00000000a80425fb
CapAmb: 0000000000000000
NoNewPrivs:     1
Seccomp:        2
`,
			want: &Status{
				Name:       "my-process",
				Ruid:       1000,
				Euid:       1001,
				Suid:       1002,
				Fsuid:      1003,
				Gid:        2000,
				Egid:       2001,
				Sgid:       2002,
				Fsgid:      2003,
				CapInh:     0x00000000a80425fb,
				CapPrm:     0x00000000a80425fb,
				CapEff:     0x00000000a80425fb,
				CapBnd:     0x00000000a80425fb,
				CapAmb:     0,
				NoNewPrivs: true,
				Seccomp:    2,
			},
		},
		{
			name: "partial input",
			content: `
Name:   partial-process
Uid:    1000
Seccomp:        1
`,
			want: &Status{
				Name:    "partial-process",
				Seccomp: 1,
			},
		},
		{
			name:    "empty input",
			content: "",
			want:    &Status{},
		},
		{
			name: "malformed lines",
			content: `
Name:   malformed-process
Uid:    1000
This is a malformed line
Gid:    2000
Another: one: with: extra: colons
`,
			want: &Status{
				Name: "malformed-process",
			},
		},
		{
			name: "extra whitespace",
			content: `
Name:      whitespace-process
    Uid:    1000    1001 1002 1003
Gid:	2000	2001 2002 2003
`,
			want: &Status{
				Name:  "whitespace-process",
				Ruid:  1000,
				Euid:  1001,
				Suid:  1002,
				Fsuid: 1003,
				Gid:   2000,
				Egid:  2001,
				Sgid:  2002,
				Fsgid: 2003,
			},
		},
		{
			name: "invalid values",
			content: `
Name:   invalid-values
Uid:    abc
Gid:    def
CapInh: not-a-hex
Seccomp: ghi
`,
			wantErr: true,
		},
		{
			name:    "nonewprivs 0",
			content: "NoNewPrivs:     0",
			want: &Status{
				NoNewPrivs: false,
			},
		},
		{
			name:    "nonewprivs 1",
			content: "NoNewPrivs:     1",
			want: &Status{
				NoNewPrivs: true,
			},
		},
		{
			name: "uid/gid multiple values",
			content: `
Uid: 1000 1001 1002 1003
Gid: 2000 2001 2002 2003
`,
			want: &Status{
				Ruid:  1000,
				Euid:  1001,
				Suid:  1002,
				Fsuid: 1003,
				Gid:   2000,
				Egid:  2001,
				Sgid:  2002,
				Fsgid: 2003,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Parse(strings.NewReader(tc.content))
			if (err != nil) != tc.wantErr {
				t.Fatalf("Parse() error = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantErr {
				if got != nil {
					t.Errorf("Parse() = %v, want nil", got)
				}
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Parse() = %v, want %v", got, tc.want)
			}
		})
	}
}
