# Starlark prometheus config generator

Generate prometheus config file from a starlark file

## Usage

- Modify `cfg/prometheus.bzl`
- Run `bazel run //cmd`

Example:

```
$ bazel run //cmd
INFO: Build options --compilation_mode, --dynamic_mode, and --test_sharding_strategy have changed, discarding analysis cache.
INFO: Analyzed target //cmd:cmd (77 packages loaded, 8320 targets configured).
INFO: Found 1 target...
Target //cmd:cmd up-to-date:
  bazel-bin/cmd/cmd_/cmd
INFO: Elapsed time: 0.875s, Critical Path: 0.27s
INFO: 3 processes: 1 internal, 2 linux-sandbox.
INFO: Build completed successfully, 3 total actions
INFO: Build completed successfully, 3 total actions
global_config:
  scrape_interval_seconds: 2
scrape_config:
- job: service_1
  metric_relabel_configs:
  - regex: (.+)
    source_labels:
    - old_1
    - old_2
    target_label: new
  metrics_path: /metrics
- job: service_2
  metric_relabel_configs:
  - regex: (.+)
    source_labels:
    - old_1
    - old_2
    target_label: new
  metrics_path: /metrics

```
