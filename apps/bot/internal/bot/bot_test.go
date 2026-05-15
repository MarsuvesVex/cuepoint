package bot

import (
	"testing"
)

func TestParseMarkerCommand(t *testing.T) {
	input, ok := parseMarkerCommand("!marker stream clip 00:00:10")
	if !ok {
		t.Fatal("expected command to parse")
	}
	if input.StreamID != "stream" || input.Label != "clip" || input.Timestamp != "00:00:10" {
		t.Fatalf("unexpected input: %+v", input)
	}
}

func TestParseMarkerCommandRejectsInvalidInput(t *testing.T) {
	if _, ok := parseMarkerCommand("!marker stream clip"); ok {
		t.Fatal("expected invalid command to be rejected")
	}
}
