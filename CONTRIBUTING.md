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
that implements the [`Provider` interface in `main.go`](./src/main.go). If your goal is to
add additional exported data, then you'll probably want to implement a new provider.

The currently existing providers are loosely focused around available API endpoints.

## Git commit messages

Commits made to this repository should follow the following guidelines look like the following:

```
!!! CLEANUP | Short message in present tense
```

- `!!! ` (optional): Three exclamation marks should be added in front of the commit message if it contains breaking changes.
- `CLEANUP |`: A commit category and separator must follow valid categories are the following:
    - `FEATURE`: New functionality was added.
    - `BUGFIX`: Broken functionality was restored.
    - `CLEANUP`: No functional changes, only non-functionals. Use this if you reformat code or edit comments.
    - `TASK`: For everything that is not one of the above.
    - `RELEASE`: Only for release versions.
