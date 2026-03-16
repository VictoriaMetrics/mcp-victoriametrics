set -e
set -o pipefail

#------

rm -rf cmd/mcp-victoriametrics/resources/vm

git clone --no-checkout --depth=1 https://github.com/VictoriaMetrics/vmdocs.git cmd/mcp-victoriametrics/resources/vm
cd cmd/mcp-victoriametrics/resources/vm

git sparse-checkout init --cone
git sparse-checkout set content
git checkout main
rm -rf ./.git
rm -rf content/anomaly-detection
rm -rf content/victorialogs
rm -rf content/victorialogs-grafana-datasource
rm -rf content/victoriatraces
rm -rf content/home
rm -rf content/helm/victoria-logs-agent
rm -rf content/helm/victoria-logs-cluster
rm -rf content/helm/victoria-logs-collector
rm -rf content/helm/victoria-logs-multilevel
rm -rf content/helm/victoria-logs-single
rm -rf content/helm/victoria-traces-cluster
rm -rf content/helm/victoria-traces-single
rm -f ./docs/Makefile ./Makefile ./LICENSE ./*.md ./*.mod ./*.sum ./*.zip ./.golangci.yml ./.wwhrd.yml ./.gitignore ./.dockerignore ./codecov.yml ./Dockerfile ./*.sh ./*.js ./*.json ./*.lock

cd -

#------
