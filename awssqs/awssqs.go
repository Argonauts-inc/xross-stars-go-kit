package awssqs

import (
	"context"

	"github.com/Argonauts-inc/xross-stars-go-kit/otel"
	"github.com/Argonauts-inc/xross-stars-go-kit/requestid"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const (
	MessageAttributeRequestID              = requestid.LogKey
	MessageAttributeDataTypeString         = "String"
	MessageSystemAttributeAWSTraceHeader   = otel.XRayTraceHeader
	MessageSystemAttributeAWSTraceDataType = "String"
)

func RequestIDFromLambdaMessage(record events.SQSMessage) string {
	attr, ok := record.MessageAttributes[MessageAttributeRequestID]
	if !ok || attr.StringValue == nil {
		return ""
	}

	return requestid.Normalize(*attr.StringValue)
}

func ContextWithLambdaMessageRequestID(ctx context.Context, record events.SQSMessage) context.Context {
	return requestid.WithContext(ctx, RequestIDFromLambdaMessage(record))
}

func ContextWithLambdaMessageTrace(ctx context.Context, record events.SQSMessage) context.Context {
	return otel.ExtractXRayTraceHeader(ctx, record.Attributes[MessageSystemAttributeAWSTraceHeader])
}

func RequestIDMessageAttributes(id string) map[string]sqstypes.MessageAttributeValue {
	if id = requestid.Normalize(id); id == "" {
		return nil
	}

	return map[string]sqstypes.MessageAttributeValue{
		MessageAttributeRequestID: {
			DataType:    aws.String(MessageAttributeDataTypeString),
			StringValue: aws.String(id),
		},
	}
}

func TraceMessageSystemAttributes(ctx context.Context) map[string]sqstypes.MessageSystemAttributeValue {
	traceHeader := otel.InjectXRayTraceHeader(ctx)
	if traceHeader == "" {
		return nil
	}

	return map[string]sqstypes.MessageSystemAttributeValue{
		MessageSystemAttributeAWSTraceHeader: {
			DataType:    aws.String(MessageSystemAttributeAWSTraceDataType),
			StringValue: aws.String(traceHeader),
		},
	}
}
