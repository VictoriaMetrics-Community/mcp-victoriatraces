set -e
set -o pipefail

#------

rm -rf cmd/mcp-victoriatraces/resources/vmsite

git clone --no-checkout --depth=1 https://github.com/VictoriaMetrics/vmsite.git cmd/mcp-victoriatraces/resources/vmsite
cd cmd/mcp-victoriatraces/resources/vmsite

git sparse-checkout init --cone
git sparse-checkout set content/blog
git checkout master
rm -rf ./.git
rm -f ./content/_index.md ./Dockerfile ./Makefile ./*.md ./*.json ./*.lock ./.gitignore

rm -rf ./content/blog/categories
rm -rf ./content/blog/tags

find content/blog/*  | grep -v tracing | grep -v traces | xargs -I {} rm -rf {}

cd -

#------
