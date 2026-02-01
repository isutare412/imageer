package tracing

import semconv "go.opentelemetry.io/otel/semconv/v1.37.0"

var (
	PeerServiceGoogleOIDC = semconv.PeerService("google-oidc")
	PeerServicePostgres   = semconv.PeerService("postgres")
	PeerServiceValkey     = semconv.PeerService("valkey")
	PeerServiceAWSSQS     = semconv.PeerService("aws-sqs")
	PeerServiceAWSS3      = semconv.PeerService("aws-s3")
	PeerServiceInternet   = semconv.PeerService("internet")
)
