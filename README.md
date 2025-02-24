# Go Template

This repository contains a Go application that can be used for rendering Go templates.

## Usage
```
go-template <-t templateFilePath> [-f valuesFilePath] [-s variableOverride] [-o outputFolder]
```

## Release

### Manual Release
To generate binaries, run:
```sh
goreleaser release --clean --skip publish
```

### Automated Release
[.goreleaser.yml](.goreleaser.yml) can be integrated with tools like Github Actions to automate build and release process.

## License

[Apache License](LICENSE)