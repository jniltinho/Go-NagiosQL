#!/usr/bin/env bash
# Monitoring summary smoke tests.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

echo "=== Monitoring Tests ==="

login admin admin123

# Summary endpoint.
assert_status "/api/v1/monitoring/summary" GET 200

# Verify key fields exist.
for field in hosts services commands timeperiods contacts hostgroups servicegroups contactgroups hosttemplates servicetemplates; do
    val=$(python3 -c "import json; d=json.load(open('/tmp/nagiosql_body.json')); print(d.get('${field}','MISSING'))" 2>/dev/null)
    if [ "$val" = "MISSING" ]; then
        echo "[FAIL] summary missing field: ${field}"
        FAIL=$((FAIL+1))
    else
        echo "[PASS] summary.${field}=${val}"
        PASS=$((PASS+1))
    fi
done

summary
