package ffmpeg

import (
	"strings"
	"testing"

	"github.com/MarsuvesVex/cuepoint/packages/stream"
)

func TestBuildSegmentCommand(t *testing.T) {
	marker := stream.Marker{ID: "marker-1", StreamID: "stream-source", Timestamp: "00:10:00"}
	job := stream.Job{ID: "job-1"}

	command := BuildSegmentCommand(marker, job)

	for _, want := range []string{"ffmpeg", "00:10:00", "stream-source", "marker-1-job-1"} {
		if !strings.Contains(command, want) {
			t.Fatalf("command %q missing %q", command, want)
		}
	}
}
