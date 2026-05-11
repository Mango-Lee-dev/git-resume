#!/usr/bin/env bash
# Start backend (Go API, :8080) and frontend (Next.js, :3000) together.
# Ctrl-C stops both.

set -e
cd "$(dirname "$0")"

LOG_DIR="logs"
mkdir -p "$LOG_DIR"
BACK_LOG="$LOG_DIR/backend.log"
FRONT_LOG="$LOG_DIR/frontend.log"

# Stop any previous instances we might have started
pkill -f "bin/git-resume serve" 2>/dev/null || true
# Kill anything holding our frontend port
lsof -ti tcp:3000 2>/dev/null | xargs kill 2>/dev/null || true

cleanup() {
  echo
  echo "Stopping..."
  # Kill immediate children (yarn wrapper, backend) and their descendants
  for pid in "${BACK_PID:-}" "${FRONT_PID:-}"; do
    [ -n "$pid" ] || continue
    pkill -TERM -P "$pid" 2>/dev/null || true
    kill "$pid" 2>/dev/null || true
  done
  # next dev forks a long-lived node child that survives killing yarn
  pkill -f "next dev" 2>/dev/null || true
  pkill -f "next-server" 2>/dev/null || true
  pkill -f "bin/git-resume serve" 2>/dev/null || true
  pkill -P $$ 2>/dev/null || true
  exit 0
}
trap cleanup INT TERM

# --- Backend ----------------------------------------------------------------
if [ ! -x bin/git-resume ]; then
  echo "Building backend..."
  go build -o bin/git-resume .
fi

if [ ! -f .env ] || ! grep -q '^CLAUDE_API_KEY=' .env; then
  echo "WARN: .env missing CLAUDE_API_KEY. Backend will start but analysis will fail until a session key is set in the UI."
fi

echo "Starting backend (logs: $BACK_LOG)"
./bin/git-resume serve > "$BACK_LOG" 2>&1 &
BACK_PID=$!

# --- Frontend ---------------------------------------------------------------
if [ ! -d web/node_modules ] || [ ! -f web/package.json ]; then
  echo "Installing frontend deps..."
  (cd web && yarn install)
fi

echo "Starting frontend (logs: $FRONT_LOG)"
(cd web && yarn dev) > "$FRONT_LOG" 2>&1 &
FRONT_PID=$!

# --- Wait for backend health -----------------------------------------------
for i in $(seq 1 20); do
  if curl -sf -o /dev/null http://localhost:8080/health; then
    break
  fi
  sleep 0.3
done

echo
echo "Backend  : http://localhost:8080  (pid $BACK_PID)"
echo "Frontend : http://localhost:3000  (pid $FRONT_PID)"
echo "Ctrl-C to stop both."
echo
echo "----- streaming logs -----"

tail -F "$BACK_LOG" "$FRONT_LOG" &
wait
