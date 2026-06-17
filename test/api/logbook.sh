#!/usr/bin/env bash
# Logbook read-only smoke tests.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

echo "=== Logbook Tests ==="

login admin admin123

# Generate some audit entries by creating and deleting a host.
assert_status "/api/v1/hosts" POST 201 \
    '{"host_name":"logbook-test-host","alias":"Logbook Test","address":"10.99.99.99","check_command":"check-host-alive","active":"1","register":"1"}'
H_ID=$(python3 -c "import json; print(json.load(open('/tmp/nagiosql_body.json')).get('id',''))")
assert_status "/api/v1/hosts/${H_ID}" DELETE 204

# List logbook.
assert_status "/api/v1/logbook" GET 200

# Filter by object_type.
assert_status "/api/v1/logbook?object_type=host" GET 200

# Filter by user.
assert_status "/api/v1/logbook?user=admin" GET 200

# Confirm at least 1 entry is present.
COUNT=$(python3 -c "import json; d=json.load(open('/tmp/nagiosql_body.json')); print(d.get('total',0))" 2>/dev/null || echo 0)
if [ "$COUNT" -gt 0 ]; then
    echo "[PASS] logbook has entries (total=${COUNT})"
    PASS=$((PASS+1))
else
    echo "[FAIL] logbook empty"
    FAIL=$((FAIL+1))
fi

summary
