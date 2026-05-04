package pipeprimitive

import "testing"

func TestPasswdPasswordOffset(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
		wantErr bool
	}{
		{
			name:    "root first line",
			content: "root:x:0:0:root:/root:/bin/bash\n",
			want:    len("root:") - 1,
		},
		{
			name:    "root after another user",
			content: "daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin\nroot:x:0:0:root:/root:/bin/bash\n",
			want:    len("daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin\nroot:") - 1,
		},
		{
			name:    "root without trailing newline",
			content: "root:x:0:0:root:/root:/bin/bash",
			want:    len("root:") - 1,
		},
		{
			name:    "missing root",
			content: "daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := passwdPasswordOffset([]byte(tt.content), "root")
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("offset = %d, want %d", got, tt.want)
			}
		})
	}
}
