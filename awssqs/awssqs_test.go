package awssqs

import (
	"context"
	"testing"

	"github.com/Argonauts-inc/xross-stars-go-kit/requestid"
	"github.com/aws/aws-lambda-go/events"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func TestRequestIDFromLambdaMessage(t *testing.T) {
	t.Parallel()

	requestID := "req-123"
	got := RequestIDFromLambdaMessage(events.SQSMessage{
		MessageAttributes: map[string]events.SQSMessageAttribute{
			MessageAttributeRequestID: {StringValue: &requestID},
		},
	})

	if got != requestID {
		t.Fatalf("RequestIDFromLambdaMessage() = %q, want %q", got, requestID)
	}
}

func TestContextWithLambdaMessageRequestID(t *testing.T) {
	t.Parallel()

	requestID := "req-123"
	ctx := ContextWithLambdaMessageRequestID(context.Background(), events.SQSMessage{
		MessageAttributes: map[string]events.SQSMessageAttribute{
			MessageAttributeRequestID: {StringValue: &requestID},
		},
	})

	got, ok := requestid.FromContext(ctx)
	if !ok {
		t.Fatal("request id missing from context")
	}
	if got != requestID {
		t.Fatalf("request id = %q, want %q", got, requestID)
	}
}

func TestRequestIDMessageAttributes(t *testing.T) {
	t.Parallel()

	attrs := RequestIDMessageAttributes(" req-123 ")
	attr, ok := attrs[MessageAttributeRequestID]
	if !ok {
		t.Fatal("request id message attribute missing")
	}
	if attr.StringValue == nil || *attr.StringValue != "req-123" {
		t.Fatalf("StringValue = %v, want req-123", attr.StringValue)
	}
}

func TestTraceMessageSystemAttributes(t *testing.T) {
	t.Parallel()

	provider := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	ctx, span := provider.Tracer("test").Start(context.Background(), "test")
	defer span.End()

	attrs := TraceMessageSystemAttributes(ctx)
	attr, ok := attrs[MessageSystemAttributeAWSTraceHeader]
	if !ok {
		t.Fatal("trace message system attribute missing")
	}
	if attr.StringValue == nil || *attr.StringValue == "" {
		t.Fatalf("trace header = %v, want non-empty", attr.StringValue)
	}

	extracted := ContextWithLambdaMessageTrace(context.Background(), events.SQSMessage{
		Attributes: map[string]string{
			MessageSystemAttributeAWSTraceHeader: *attr.StringValue,
		},
	})
	if !trace.SpanContextFromContext(extracted).IsValid() {
		t.Fatal("extracted span context is invalid")
	}
}
