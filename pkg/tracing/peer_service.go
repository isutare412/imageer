package tracing

import semconv "go.opentelemetry.io/otel/semconv/v1.40.0"

var (
	PeerServiceGoogleOIDC = semconv.ServicePeerName("google-oidc")
	PeerServicePostgres   = semconv.ServicePeerName("postgres")
	PeerServiceValkey     = semconv.ServicePeerName("valkey")
	PeerServiceAWSSQS     = semconv.ServicePeerName("aws-sqs")
	PeerServiceAWSS3      = semconv.ServicePeerName("aws-s3")
	PeerServiceInternet   = semconv.ServicePeerName("internet")
)
