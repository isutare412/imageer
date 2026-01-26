package webv2

//go:generate go tool oapi-codegen -config generate-server.yaml openapi.yaml
//go:generate go tool oapi-codegen -config generate-client.yaml openapi.yaml
//go:generate go tool mockgen -package gateway -source=../../../pkg/gateway/client.gen.go -destination=../../../pkg/gateway/client.gen_mock.go
