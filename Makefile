all: tidy proto gazelle_deps gazelle_fix
	bazel run //cmd


gazelle_deps: go.mod tidy
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies

gazelle_fix: gazelle_deps
	bazel run //:gazelle -- fix

tidy:
	go mod tidy

proto: proto/github/com/michaelboulton/starlark_config_example/v1/prometheus/prometheus.proto
	mkdir -p v1/prometheus
	protoc --proto_path=proto/github/com/michaelboulton/starlark_config_example/v1/ --go_out=v1/ --go_opt=paths=source_relative prometheus/prometheus.proto
