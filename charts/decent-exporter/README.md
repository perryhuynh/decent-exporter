# decent-exporter

Prometheus exporter for [Reaprime](https://github.com/tadelv/reaprime)-powered Decent espresso tablets.

## Installing

```sh
helm install decent-exporter oci://ghcr.io/perryhuynh/charts/decent-exporter \
  --set config.reaprimeUrl=http://192.168.50.49:8080
```

## Monitoring

The chart renders a `ServiceMonitor` by default for kube-prometheus-stack scraping. Set `monitoring.serviceMonitor.enabled=false` if the Prometheus Operator CRDs are not installed.

Set `monitoring.dashboards.enabled=true` to render Grafana dashboard ConfigMaps for the sidecar pattern.

Set `monitoring.dashboards.grafanaOperator.enabled=true` and provide `monitoring.dashboards.grafanaOperator.matchLabels` for grafana-operator `GrafanaDashboard` resources.

## Values

| Key | Type | Default | Description |
| --- | --- | --- | --- |
| `image.repository` | string | `ghcr.io/perryhuynh/reaprime-exporter` | Image repository. |
| `image.tag` | string | `""` | Overrides the image tag; defaults to chart appVersion. |
| `image.digest` | string | `""` | Pin image by digest. Set by release workflow for published charts. |
| `config.port` | int | `8080` | Exporter HTTP listen port. |
| `config.reaprimeUrl` | string | `http://192.168.50.49:8080` | Base URL for the Reaprime tablet webserver. |
| `monitoring.serviceMonitor.enabled` | bool | `true` | Create a Prometheus Operator ServiceMonitor. |
| `monitoring.dashboards.enabled` | bool | `false` | Render Grafana dashboard ConfigMaps. |
| `monitoring.dashboards.grafanaOperator.enabled` | bool | `false` | Render GrafanaDashboard CRs for grafana-operator. |
