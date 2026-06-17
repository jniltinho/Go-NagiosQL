#!/usr/bin/env bash
# Host CRUD smoke tests.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

echo "=== Host Tests ==="

login admin admin123

# Create host.
echo "[STEP] create host web-test-01"
assert_status "/api/v1/hosts" POST 201 \
    '{"host_name":"web-test-01","alias":"Web Test 01","address":"10.99.99.1","check_command":"check-host-alive","active":"1","register":"1"}'

HOST_ID=$(python3 -c "import json; d=json.load(open('/tmp/nagiosql_body.json')); print(d.get('id',''))")

# List hosts — web-test-01 must appear.
assert_status "/api/v1/hosts" GET 200
if python3 -c "import json,sys; d=json.load(open('/tmp/nagiosql_body.json')); sys.exit(0 if any(h['host_name']=='web-test-01' for h in d.get('data',[])) else 1)" 2>/dev/null; then
    echo "[PASS] web-test-01 in list"
    PASS=$((PASS+1))
else
    echo "[FAIL] web-test-01 not found in list"
    FAIL=$((FAIL+1))
fi

# Get by ID.
assert_status "/api/v1/hosts/${HOST_ID}" GET 200
assert_field "host_name" "web-test-01"

# Update alias.
assert_status "/api/v1/hosts/${HOST_ID}" PUT 200 \
    "{\"host_name\":\"web-test-01\",\"alias\":\"Updated Alias\",\"address\":\"10.99.99.1\",\"check_command\":\"check-host-alive\",\"active\":\"1\",\"register\":\"1\"}"

# Delete.
assert_status "/api/v1/hosts/${HOST_ID}" DELETE 204

# Get deleted — must 404.
assert_status "/api/v1/hosts/${HOST_ID}" GET 404

summary
