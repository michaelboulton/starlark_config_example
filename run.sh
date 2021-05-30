set -ex

mkdir -p v1/prometheus
protoc --proto_path=proto/github/com/michaelboulton/starlark_config_example/v1/ --go_out=v1/ --go_opt=paths=source_relative prometheus/prometheus.proto

go mod tidy || true
bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies
bazel run //:gazelle -- fix
bazel run //cmd
