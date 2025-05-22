CLIENT_COUNT=5
CLIENT_SCRIPT="${PROJECT_ROOT}/python-client/client.py"

for i in $(seq 1 $CLIENT_COUNT); do
  USER_ID="user${i}"
  echo "Starting Python Client ${USER_ID}..."
  nohup python3 "$CLIENT_SCRIPT" > "${LOG_DIR}/${USER_ID}.log" 2>&1 <<EOF &
${USER_ID}
0
EOF
  # envia USER_ID na entrada e sai em seguida 
done

echo "All clients started (${CLIENT_COUNT}). Logs in ${LOG_DIR}."
