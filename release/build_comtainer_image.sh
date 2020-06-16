parentdir="$(dirname "$(pwd)")"

docker build --tag azure-devops-exporter:latest -f ${parentdir}/Dockerfile ${parentdir}
