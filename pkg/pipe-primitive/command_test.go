package pipeprimitive

import "testing"

func TestCommandExploitSubcommandsUseShortNames(t *testing.T) {
	cmd := Command(&recordingPrimitive{}, nil, "recording primitive")

	topLevel := map[string]bool{}
	for _, sub := range cmd.Commands {
		topLevel[sub.Name] = true
	}
	for _, want := range []string{"privilege-escalate", "escape", "image-pollution", "clean"} {
		if !topLevel[want] {
			t.Fatalf("top-level exploit subcommand %q missing from %+v", want, topLevel)
		}
	}
	for _, unwanted := range []string{"recording-privilege-escalate", "recording-escape", "recording-image-pollution"} {
		if topLevel[unwanted] {
			t.Fatalf("top-level exploit subcommand %q should use a short name", unwanted)
		}
	}

	var escapeCmdFound bool
	for _, sub := range cmd.Commands {
		if sub.Name != "escape" {
			continue
		}
		escapeCmdFound = true

		got := map[string]bool{}
		for _, escapeSub := range sub.Commands {
			got[escapeSub.Name] = true
		}
		for _, want := range []string{"image", "exec", "restart"} {
			if !got[want] {
				t.Fatalf("escape subcommand %q missing from %+v", want, got)
			}
		}
	}
	if !escapeCmdFound {
		t.Fatal("escape command missing")
	}
}
