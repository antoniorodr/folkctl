package cmd

import "testing"

func TestDetectClipboardCmdForOS(t *testing.T) {
	tests := []struct {
		name string
		goos string
		want string
	}{
		{"macOS", "darwin", "pbcopy"},
		{"Linux", "linux", "xclip -selection clipboard"},
		{"Windows", "windows", "clip"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectClipboardCmdForOS(tt.goos)
			if got != tt.want {
				t.Errorf("detectClipboardCmdForOS(%q) = %q, want %q", tt.goos, got, tt.want)
			}
		})
	}
}
