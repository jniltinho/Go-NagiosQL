#!/usr/bin/env bash
# Command CRUD smoke tests.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

echo "=== Command Tests ==="

login admin admin123

# Create check command (type=0).
assert_status "/api/v1/commands" POST 201 \
    '{"command_name":"smoke-check-ping","command_line":"$USER1$/check_ping -H $HOSTADDRESS$ -w 100,20% -c 500,60%","command_type":0,"active":"1","register":"1"}'
CHECK_ID=$(python3 -c "import json; print(json.load(open('/tmp/nagiosql_body.json')).get('id',''))")

# Create notify command (type=1).
assert_status "/api/v1/commands" POST 201 \
    '{"command_name":"smoke-notify-email","command_line":"/usr/bin/mail -s $HOSTNAME$ $CONTACTEMAIL$","command_type":1,"active":"1","register":"1"}'
NOTIFY_ID=$(python3 -c "import json; print(json.load(open('/tmp/nagiosql_body.json')).get('id',''))")

# List check commands only.
assert_status "/api/v1/commands?type=check" GET 200
if python3 -c "import json,sys; d=json.load(open('/tmp/nagiosql_body.json')); sys.exit(0 if any(c['command_name']=='smoke-check-ping' for c in d.get('data',[])) else 1)" 2>/dev/null; then
    echo "[PASS] smoke-check-ping in check list"
    PASS=$((PASS+1))
else
    echo "[FAIL] smoke-check-ping not found in check list"
    FAIL=$((FAIL+1))
fi

# List notify commands only.
assert_status "/api/v1/commands?type=notify" GET 200

# Delete both.
assert_status "/api/v1/commands/${CHECK_ID}" DELETE 204
assert_status "/api/v1/commands/${NOTIFY_ID}" DELETE 204

summary
