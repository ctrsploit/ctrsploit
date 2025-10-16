package kernel

import "testing"

func Test_getCompiledDate(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
		wantErr  bool
	}{
		{
			name: "Ubuntu style",
			content: `Linux version 6.5.0-14-generic (buildd@lcy02-amd64)
			(gcc (Ubuntu 12.3.0-1ubuntu1~23.04) 12.3.0, GNU ld 2.40)
			#14-Ubuntu SMP Mon Jun 5 14:18:34 UTC 2023`,
			expected: "2023-06-05",
		},
		{
			name: "RedHat RFC1123Z style",
			content: `Linux version 5.19.0-46-generic (builder@rhel)
			(gcc version 11.3.0) #1 SMP Tue, 13 Feb 2024 15:21:03 +0000`,
			expected: "2024-02-13",
		},
		{
			name: "Kali style (simple date)",
			content: `Linux version 9.99.9-amd64 (devel@kali.org)
			(x86_64-linux-gnu-gcc-14 (Debian 99.9.9-9) 99.9.9, GNU ld 9.99.9)
			#1 SMP PREEMPT_DYNAMIC Kali 6.11.2-1kali1 (2099-11-22`,
			expected: "2099-11-22",
		},
		{
			name:    "Unknown style (should fail)",
			content: `Linux version 6.11.2 (no date info here)`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiledDate, err := getCompiledDate(tt.content)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			got := compiledDate.Format("2006-01-02")
			if got != tt.expected {
				t.Errorf("expected date %v, got %v", tt.expected, got)
			}
		})
	}
}
