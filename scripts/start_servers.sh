# Inicia 3 instâncias do Go Server e 1 instância do Java Coordinator

set -e
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="${SCRIPT_DIR}/.."
LOG_DIR="${PROJECT_ROOT}/logs"

mkdir -p "${LOG_DIR}"

# Servidores Go escutam PULL em porta 5560, REP em 5555 e PUB em 5556

SERVERS=("server0" "server1" "server2")
BASE_PULL_PORT=5560
BASE_REP_PORT=5555
BASE_PUB_PORT=5556

# Funcão para construir PEER_ADDRS para um servidor
build_peers() {
  local self="$1"
  local addrs=()
  for s in "${SERVERS[@]}"; do
    if [[ "$s" != "$self" ]]; then
      addrs+=("tcp://localhost:${BASE_PULL_PORT}")
    fi
  done
  IFS=","; echo "${addrs[*]}"; IFS=$' \t\n'
}

# Inicia instâncias Go Server
for srv in "${SERVERS[@]}"; do
  PEERS=$(build_peers "$srv")
  echo "Starting Go Server ${srv} with PEER_ADDRS=${PEERS}..."
  SERVER_ID=${srv} PEER_ADDRS=${PEERS} \
    nohup "${PROJECT_ROOT}/go-server/server" \
    > "${LOG_DIR}/${srv}.log" 2>&1 &
done

# Delay para garantir que Go servers subam antes do coordinator
sleep 2

# Inicia Java Coordinator
COORD_ID="coordinator"
# Lista de servidores para coord
COORD_LIST=$(IFS=","; echo "${SERVERS[*]}")

echo "Starting Java Coordinator ${COORD_ID} for [${COORD_LIST}]..."
nohup java -cp "${PROJECT_ROOT}/java-coordinator/*:your_zmq_and_jackson_jars" \
  Coordinator "${COORD_ID}" "${COORD_LIST}" \
  > "${LOG_DIR}/coordinator.log" 2>&1 &

echo "All servers started. Logs in ${LOG_DIR}."
