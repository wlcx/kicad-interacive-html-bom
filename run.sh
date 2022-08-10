#!/usr/bin/env bash

generate() {
    mkdir gen
    LD_LIBRARY_PATH="/usr/lib/kicad-nightly/lib/x86_64-linux-gnu/" PYTHONPATH="/usr/lib/kicad-nightly/lib/python3/dist-packages" xvfb-run -a python3 /opt/InteractiveHtmlBom-master/InteractiveHtmlBom/generate_interactive_bom.py --no-browser --name-format "%f-iBOM" --dest-dir ./gen "$1"
}

for f in *.kicad_pcb; do
    generate "$f"
done

for f in *.kicad_sch; do
    eeschema_do export -a "$f" gen/
done

mkdir out
/opt/kicadsitegenerator -projectName $GITHUB_REPOSITORY -projectVersion "${GITHUB_SHA:0:8}" -out out *.md gen/*.pdf gen/*.html
