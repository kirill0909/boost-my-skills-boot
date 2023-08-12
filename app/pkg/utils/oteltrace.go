package utils

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

const (
	traceIDHeader string = "X-Trace-ID"
)

func StartFiberTrace(c *fiber.Ctx, spanName string) (context.Context, trace.Span) {
	ctx, span := otel.Tracer("").Start(c.Context(), spanName)
	c.Set(traceIDHeader, span.SpanContext().SpanID().String())
	if clientID, ok := c.Locals("clientID").(int); ok {
		span.SetAttributes(attribute.Int("clientID", clientID))
	}
	if clientID, ok := c.Locals("APIID").(int); ok {
		span.SetAttributes(attribute.Int("APIID", clientID))
	}
	if userID, ok := c.Locals("userID").(int); ok {
		span.SetAttributes(attribute.Int("userID", userID))
	}
	return ctx, span
}

func StartGrpcTrace(ctx context.Context, spanName string) (context.Context, trace.Span) {
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	ctx = metadata.AppendToOutgoingContext(ctx, traceIDHeader, span.SpanContext().TraceID().String())
	return ctx, span
}
