# decent-exporter

Prometheus exporter for [Reaprime](https://github.com/tadelv/reaprime)-powered Decent espresso tablets.

The exporter connects to [Reaprime](https://github.com/tadelv/reaprime)'s local WebSocket API, keeps live snapshots in memory, and exposes stable Prometheus metrics on `/metrics`.

## Run

```sh
DECENT_EXPORTER_URL=http://192.168.50.49:8080 go run ./cmd/decent-exporter
```

```sh
curl http://127.0.0.1:8080/metrics
```

## Configuration

| Environment variable | Default | Description |
| --- | --- | --- |
| `DECENT_EXPORTER_LISTEN_ADDRESS` | `:8080` | Exporter listen address. |
| `DECENT_EXPORTER_URL` | `http://127.0.0.1:8080` | Base URL for the Reaprime tablet webserver. |
| `DECENT_EXPORTER_LOG_LEVEL` | `info` | `info` or `debug`. |
| `DECENT_EXPORTER_READY_MAX_AGE` | `30s` | Maximum age of the machine stream before `/readyz` fails. |
| `DECENT_EXPORTER_RECONNECT_MIN` | `1s` | Initial stream reconnect delay. |
| `DECENT_EXPORTER_RECONNECT_MAX` | `30s` | Maximum stream reconnect delay. |

## Metrics scope

The exporter intentionally avoids labels containing serial numbers, device IDs, bean names, profile names, notes, logs, or raw BLE payloads.

Pressure, flow, and weight change sub-second during a shot, but Prometheus scrapes every 15–60s — at that resolution instantaneous gauges of those signals are aliased noise. Instead of exporting them raw, the exporter watches every Reaprime machine sample (~1.7/s) and emits per-shot summaries that persist between scrapes:

- `decent_shots_total` — completed espresso shots since start.
- `decent_shot_duration_seconds`, `decent_shot_peak_pressure_bar`, `decent_shot_peak_flow_ml_per_second`, `decent_shot_average_flow_ml_per_second` — stats for the last completed shot.

A shot is the span the machine spends in the `espresso` state; episodes shorter than 3s (flushes, aborts) are discarded. The scale is not exported at all.
