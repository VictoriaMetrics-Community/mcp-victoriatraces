set -e
set -o pipefail

#------

rm -rf cmd/mcp-victoriatraces/resources/vm

git clone --no-checkout --depth=1 https://github.com/VictoriaMetrics/VictoriaTraces.git cmd/mcp-victoriatraces/resources/vm
cd cmd/mcp-victoriatraces/resources/vm

git sparse-checkout init --cone
git sparse-checkout set docs
git checkout master
rm -rf ./.git
rm -f ./docs/Makefile ./Makefile ./LICENSE ./*.md ./*.mod ./*.sum ./*.zip ./.golangci.yml ./.wwhrd.yml ./.gitignore ./.dockerignore ./codecov.yml

cd -

#------
