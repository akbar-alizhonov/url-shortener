package logger

import "context"

const RequestID = "request_id"

func WithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, "request_id", reqID)
}

func RequestIDFromContext(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestID).(string); ok {
		return reqID
	}
	return ""
}
