pb = proto.package("github.com.michaelboulton.starlark_config_example.v1")

relabel = pb.RelabelConfig(
    source_labels = [
        "old_1",
        "old_2",
    ],
    target_label = "new",
    regex = "(.+)",
)

def generic_service(job):
    return pb.ScrapeConfig(
        job = job,
        metrics_path = "/metrics",
        metric_relabel_configs = [relabel],
    )

def main(ctx):
    global_cfg = pb.GlobalConfig(scrape_interval_seconds = 2)

    scrapes_cfg = [
        generic_service("service_1"),
        generic_service("service_2"),
    ]

    cfg = pb.WholeConfig(
        global_config = global_cfg,
        scrape_config = scrapes_cfg,
    )

    return [cfg]
