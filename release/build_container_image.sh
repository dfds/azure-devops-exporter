parentdir="$(dirname "$(pwd)")"

docker build --tag ded/azure-devops-exporter:latest -f ${parentdir}/Dockerfile ${parentdir}
