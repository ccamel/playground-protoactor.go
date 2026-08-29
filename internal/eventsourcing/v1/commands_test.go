package eventsourcingv1

import (
	"testing"
	"time"

	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCommandStatusMetadataRoundTrip(t *testing.T) {
	t.Parallel()

	processedAt := timestamppb.New(time.Date(2026, time.August, 29, 16, 9, 0, 0, time.UTC))
	status := &CommandStatus{
		Code:        code.Code_OK,
		Message:     "command processed",
		ProcessedAt: processedAt,
		CommandId:   "command-123",
	}

	payload, err := proto.Marshal(status)
	if err != nil {
		t.Fatalf("marshal status: %v", err)
	}

	var got CommandStatus
	if err := proto.Unmarshal(payload, &got); err != nil {
		t.Fatalf("unmarshal status: %v", err)
	}

	if !proto.Equal(status, &got) {
		t.Errorf("round trip status = %v, want %v", &got, status)
	}
}
