# KiCAD Site Generator
This project can generate a static site from a KiCAD project repo, rendering markdown
documents, PDF schematics and [Interactive HTML BOMs][ibom] and wrapping them in a nice
UI.
It is intended to run as a Github action. Support for other use cases may come at some
point (contributions welcome!). 

## Usage
Example from a Github actions workflow yaml file:
```yaml
jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
        name: Checkout repo
      - uses: wlcx/kicad-interactive-html-bom@main
        name: Run kicad interactive HTML BOM
      - uses: actions/upload-artifact@v2
        name: Upload artifacts
        with:
          name: Cool site
          path: out/*
```

## Testing locally
You can test the workflow locally by building the container and running it with a
mounted directory:

- (In local clone) `docker build -t kicadsitegen .`
- `cd` to your kicad project's directory
- `docker run --rm -it -v ``pwd``:/opt/project kicadsitegen`

[ibom]: https://github.com/openscopeproject/InteractiveHtmlBom
