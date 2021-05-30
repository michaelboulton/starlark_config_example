set -e

go mod tidy || true
bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies
bazel run //:gazelle -- fix
bazel run //cmd
