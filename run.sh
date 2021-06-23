#!/usr/bin/env bash

generate() {
    mkdir gen
    LD_LIBRARY_PATH="/usr/lib/kicad-nightly/lib/x86_64-linux-gnu/" PYTHONPATH="/usr/lib/kicad-nightly/lib/python3/dist-packages" xvfb-run -a python3 /opt/InteractiveHtmlBom-master/InteractiveHtmlBom/generate_interactive_bom.py --no-browser --dest-dir ./gen "$1"
}

cat <<EOT >  /opt/InteractiveHtmlBom-master/InteractiveHtmlBom/web/userheader.html
<div style="width:100%; height: 20px;">
    Generated from git: ${GIT_HASH:-unknown}
</div>
EOT

for f in *.kicad_pcb; do
    generate "$f"
done
