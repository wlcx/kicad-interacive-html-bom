#!/usr/bin/env bash

generate() {
    mkdir gen
    LD_LIBRARY_PATH="/usr/lib/kicad-nightly/lib/x86_64-linux-gnu/" PYTHONPATH="/usr/lib/kicad-nightly/lib/python3/dist-packages" xvfb-run -a python3 /opt/InteractiveHtmlBom-master/InteractiveHtmlBom/generate_interactive_bom.py --no-browser --name-format "%f" --dest-dir ./gen "$1"
}

cat <<EOT >  /opt/InteractiveHtmlBom-master/InteractiveHtmlBom/web/user.js
document.addEventListener('DOMContentLoaded', function() {
  for (ele of document.querySelectorAll(".fileinfo .title")) {
    ele.style.fontSize = "16px";
  }
  row = document.createElement("tr");
  genDate = document.createElement("td");
  genDate.innerHTML = "Generated: $(date --rfc-3339=seconds)";
  row.appendChild(genDate);
  gitInfo = document.createElement("td");
  gitInfo.innerHTML = "git: <a href=\"https://github.com/${GITHUB_REPOSITORY}/commit/${GITHUB_SHA}\">${GITHUB_SHA:0:8-unknown}</a>";
  row.appendChild(gitInfo);
  document.querySelector(".fileinfo>tbody").appendChild(row);
}, false);
EOT

for f in *.kicad_pcb; do
    generate "$f"
done
