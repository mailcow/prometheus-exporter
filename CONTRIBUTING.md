# Contribution Guidelines

* [CHANGELOG](./CHANGELOG.md)
* [LICENSE](./LICENSE)
* [README](./README.md)
* [CONTRIBUTING](./CONTRIBUTING.md)

Contributions are always welcome - just open an issue or PR.
Please be aware that even though it may not seem like it there are always other humans on the other end of a conversation.

**Please treat others with respect and kindness.**

If you have general questions or suggestions for improvement, open an issue or
[send me a mail](mailto:mailcow-exporter@j6s.dev)

## Adding a new Provider

The basic architecture of the exporter is centered around a provider. A provider is a struct
that implements the [`Provider` interface in `lib/provider/provider.go`](./lib/provider/provider.go). If your goal is to
add additional exported data, then you'll probably want to implement a new provider.

The currently existing providers are loosely focused around available API endpoints.
