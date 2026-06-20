# Changelog

## [1.0.0](https://github.com/perryhuynh/decent-exporter/compare/0.2.1...1.0.0) (2026-06-20)


### ⚠ BREAKING CHANGES

* env var renamed to DECENT_EXPORTER_URL and metrics renamed from decent_reaprime_stream_* to decent_stream_*; update deployments, dashboards, and alerts.

### Features

* **metrics:** replace raw pressure/flow/scale with per-shot summaries ([#9](https://github.com/perryhuynh/decent-exporter/issues/9)) ([278580b](https://github.com/perryhuynh/decent-exporter/commit/278580be39721ea1ac899ceb44da2424cf4d59d0))


### Code Refactoring

* rename reaprime identifiers to decent ([#11](https://github.com/perryhuynh/decent-exporter/issues/11)) ([f2334a1](https://github.com/perryhuynh/decent-exporter/commit/f2334a182cf63d8b0bde03125250c13e696a2255))

## [0.2.1](https://github.com/perryhuynh/reaprime-exporter/compare/0.2.0...0.2.1) (2026-06-19)


### Bug Fixes

* **chart:** point image.repository at reaprime-exporter ([#7](https://github.com/perryhuynh/reaprime-exporter/issues/7)) ([bbab353](https://github.com/perryhuynh/reaprime-exporter/commit/bbab3531047292bfd79b4100c86c91f644a10045)), closes [#6](https://github.com/perryhuynh/reaprime-exporter/issues/6)

## [0.2.0](https://github.com/perryhuynh/reaprime-exporter/compare/0.1.0...0.2.0) (2026-06-19)


### Features

* bootstrap reaprime exporter ([ab08f89](https://github.com/perryhuynh/reaprime-exporter/commit/ab08f8969cfd8364f1a4ee1465305e94eae988b5))
* **container:** update image docker/dockerfile (1.18 → 1.25) ([2660cc7](https://github.com/perryhuynh/reaprime-exporter/commit/2660cc7abc3b408e8f85f64b911066bb31832beb))
* **container:** update image docker/dockerfile (1.18 → 1.25) ([77f3060](https://github.com/perryhuynh/reaprime-exporter/commit/77f30600b0fd97da998ee60c8c6039cfd513a861))


### Bug Fixes

* add mise lockfile ([dc68158](https://github.com/perryhuynh/reaprime-exporter/commit/dc6815845d468b1796a6e7c537519596d6950c22))
* **ci:** emit bare semver release tags ([#4](https://github.com/perryhuynh/reaprime-exporter/issues/4)) ([99059b6](https://github.com/perryhuynh/reaprime-exporter/commit/99059b6d78263c11687373e75c1ca68cbb6e9bd1))
* complete locked mise ci setup ([99c3e37](https://github.com/perryhuynh/reaprime-exporter/commit/99c3e3767c4fa797af2e09e2f5c4db95b977a448))
* disable mise hooks in ci ([6996ba1](https://github.com/perryhuynh/reaprime-exporter/commit/6996ba1fa6eaf39fa6bcaee3f7c776230ad860c2))

## [0.1.0](https://github.com/perryhuynh/reaprime-exporter/compare/decent-exporter-v0.0.1...decent-exporter-v0.1.0) (2026-06-19)


### Features

* bootstrap reaprime exporter ([ab08f89](https://github.com/perryhuynh/reaprime-exporter/commit/ab08f8969cfd8364f1a4ee1465305e94eae988b5))
* **container:** update image docker/dockerfile (1.18 → 1.25) ([2660cc7](https://github.com/perryhuynh/reaprime-exporter/commit/2660cc7abc3b408e8f85f64b911066bb31832beb))
* **container:** update image docker/dockerfile (1.18 → 1.25) ([77f3060](https://github.com/perryhuynh/reaprime-exporter/commit/77f30600b0fd97da998ee60c8c6039cfd513a861))


### Bug Fixes

* add mise lockfile ([dc68158](https://github.com/perryhuynh/reaprime-exporter/commit/dc6815845d468b1796a6e7c537519596d6950c22))
* complete locked mise ci setup ([99c3e37](https://github.com/perryhuynh/reaprime-exporter/commit/99c3e3767c4fa797af2e09e2f5c4db95b977a448))
* disable mise hooks in ci ([6996ba1](https://github.com/perryhuynh/reaprime-exporter/commit/6996ba1fa6eaf39fa6bcaee3f7c776230ad860c2))
