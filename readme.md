# KiCAD Interactive HTML BOM Github Action
This repo is a Github action for the KiCAD [Interactive HTML BOM plugin][ibom]. It will generate an HTML BOM file for each `kicad_pcb` file in the root of the repository.

## Usage
Example from a workflow yaml file:
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
          name: Interactive BOMs
          path: gen/*.html
```

## Tweaks
We inject some [custom js][ibomcustom] into the HTML to add a git hash and generation timestamp to the header of the page. Convenient eh?

[ibom]: https://github.com/openscopeproject/InteractiveHtmlBom
[ibomcustom]: https://github.com/openscopeproject/InteractiveHtmlBom/wiki/Customization
