# Changelog

* [CHANGELOG](./CHANGELOG.md)
* [LICENSE](./LICENSE)
* [README](./README.md)
* [CONTRIBUTING](./CONTRIBUTING.md)

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.1] - 2025-03-12
### Fixed
* Fixed accidentally publishing windows/arm64 binary as windows/amd64

## [2.0.0] - 2025-03-12

With this release, `mailcow-exporter` becomes part of the official mailcow organization, which
means some things are renamed. Also, a bit of modernization was done on the codebase.

### Breaking Changes
* In order to simplify usage and prevent security issues, configuration can no longer
  be provided as URL parameters - it must be set either through environment variables
  or CLI flags when starting the exporter.
* The exporter now uses a token to secure access by default. See the README for the new
  recommended setup.
* The following CLI flags have been renamed:
  * `--defaultHost` is now `--host`
  * `--apikey` is now `--api-key`
* Docker images have moved from dockerhub to `ghcr.io/mailcow/prometheus-exporter`


When previously using the following prometheus config:

```yaml
scrape_configs:
  - job_name: 'mailcow'
    static_configs:
      - targets: ['mailcow-exporter-hostname:9099']
    params:
      host: ['mail.mycompany.com']
      apiKey: ['abc123']
```

You must now start the exporter either with `--host=mail.mycompany.com` and `--api-key=abc123`
or `MAILCOW_EXPORTER_HOST=mail.mycompany.com` and `MAILCOW_EXPORTER_API_KEY=abc123`. The
prometheus config can be reduced to the following:

```yaml
scrape_configs:
  - job_name: 'mailcow'
    static_configs:
      - targets: ['mailcow-exporter-hostname:9099']
    params:
      token: ['abc123']   # Please read the section about token authentication in the README
```


### Added
* `--scheme` can now be provided via CLI flag

### Docker images
- `ghcr.io/mailcow/prometheus-exporter:2`
- `ghcr.io/mailcow/prometheus-exporter:2.0`
- `ghcr.io/mailcow/prometheus-exporter:2.0.0`

## [1.4.0] - 2023-12-07
### Added
* Command line options `-defaultHost`, `-apikey`, `-listen` can now be set by environment variables
  `MAILCOW_EXPORTER_HOST`, `MAILCOW_EXPORTER_API_KEY`, `MAILCOW_EXPORTER_LISTEN`

## [1.3.1] - 2021-07-14
### Fixed
* A recent version of mailcow changed the type of a property from a string to an int. This
  release adds support for newer mailcow versions returning an int while preserving functionality
  on older mailcow versions that return a string.

## [1.3.0] - 2021-05-18
### Added
* New `scheme` option to allow API requests via http. (Thank you [maximbaz](https://github.com/maximbaz))

## [1.2.0] - 2020-09-06
### Added
* New rspamd metrics. Requires an up-to-date mailcow version as it uses a brand new API endpoint.

## [1.1.3] - 2020-09-06
### Fixed
* Errors in single providers will no longer translate to errors in the whole exporter.
  Instead, the new `mailcow_exporter_success` and `mailcow_api_success` metrics will then
  be set to 0. This is done to make the exporter provide metrics, even if parts of it fail.

## [1.1.2] - 2020-09-05
### Fixed
* non-200 API responses (such as authorization errors) no longer throw an obscure JSON
  marshalling error, but a more helpful message. In general, API error messages contain
  a lot more information now, which makes finding issues in ones setup easier.

## [1.1.1] - 2020-09-05
### Fixed
* `mailcow_container_start` accidentally reported a static value

## [1.1.0] - 2020-09-05
### Added
* Meta information about the mailcow API requests
* Container information
* Help texts
