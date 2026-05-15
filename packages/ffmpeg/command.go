package ffmpeg

import (
	"fmt"

	"github.com/MarsuvesVex/cuepoint/packages/stream"
)

func BuildSegmentCommand(marker stream.Marker, job stream.Job) string {
	return fmt.Sprintf(
		"ffmpeg -ss %s -i %q -c copy -f segment -segment_time 30 %q",
		marker.Timestamp,
		marker.StreamID,
		fmt.Sprintf("%s-%s-%%03d.ts", marker.ID, job.ID),
	)
}
