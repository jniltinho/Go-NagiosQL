#!/usr/bin/env bash
# Master smoke test runner — executes all test scripts in order.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

export BASE_URL="${BASE_URL:-http://localhost:8081}"
TOTAL_FAIL=0

run_script() {
    local script="$1"
    echo ""
    echo "########################################"
    echo "# Running: ${script}"
    echo "########################################"
    if bash "${SCRIPT_DIR}/${script}"; then
        echo "  [OK] ${script} passed"
    else
        echo "  [ERROR] ${script} FAILED"
        TOTAL_FAIL=$((TOTAL_FAIL+1))
    fi
}

# Health check first.
echo "[STEP] health check"
if curl -sf "${BASE_URL}/healthz" > /dev/null; then
    echo "[PASS] /healthz OK"
else
    echo "[FAIL] server not reachable at ${BASE_URL}"
    exit 1
fi

run_script auth.sh
run_script commands.sh
run_script hosts.sh
run_script services.sh
run_script timeperiods.sh
run_script contacts.sh
run_script groups.sh
run_script users.sh
run_script logbook.sh
run_script monitoring.sh
run_script import.sh

echo ""
echo "========================================"
echo "Smoke tests complete. Failed scripts: ${TOTAL_FAIL}"
echo "========================================"

if [ "$TOTAL_FAIL" -gt 0 ]; then
    exit 1
fi
