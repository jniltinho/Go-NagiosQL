#!/usr/bin/env bash
# Import endpoint smoke tests.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

echo "=== Import Tests ==="

login admin admin123

# Create a minimal Nagios config file to import.
TMPFILE=$(mktemp /tmp/nagiosql-import-XXXXXX.cfg)
cat > "$TMPFILE" << 'EOF'
define host {
    host_name               smoke-import-host
    alias                   Smoke Import Host
    address                 10.88.0.1
    check_command           check-host-alive
    max_check_attempts      3
    check_interval          5
    retry_interval          1
    check_period            24x7
    notification_interval   60
    notification_period     24x7
    contact_groups          admins
    active_checks_enabled   1
    passive_checks_enabled  1
    notifications_enabled   1
    register                1
}

define service {
    host_name               smoke-import-host
    service_description     PING
    check_command           check_ping!100,20%!500,60%
    max_check_attempts      3
    check_interval          5
    retry_interval          1
    check_period            24x7
    notification_interval   60
    notification_period     24x7
    contact_groups          admins
    active_checks_enabled   1
    passive_checks_enabled  1
    notifications_enabled   1
    register                1
}

define command {
    command_name    smoke-import-check
    command_line    $USER1$/check_ping -H $HOSTADDRESS$ -w $ARG1$ -c $ARG2$
}
EOF

# POST /api/v1/import with the config file content.
CONTENT=$(cat "$TMPFILE")
assert_status "/api/v1/import" POST 200 \
    "{\"content\":$(python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))' < "$TMPFILE"),\"config_id\":0,\"overwrite\":false}"

# Verify stats in response.
INSERTED=$(python3 -c "import json; d=json.load(open('/tmp/nagiosql_body.json')); print(d.get('inserted',0))" 2>/dev/null || echo 0)
echo "[INFO] import result: inserted=${INSERTED}"
if [ "$INSERTED" -gt 0 ]; then
    echo "[PASS] import created ${INSERTED} objects"
    PASS=$((PASS+1))
else
    echo "[FAIL] import inserted 0 objects"
    FAIL=$((FAIL+1))
fi

# Second import with overwrite=false → should skip (not insert duplicates).
assert_status "/api/v1/import" POST 200 \
    "{\"content\":$(python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))' < "$TMPFILE"),\"config_id\":0,\"overwrite\":false}"

SKIPPED=$(python3 -c "import json; d=json.load(open('/tmp/nagiosql_body.json')); print(d.get('skipped',0))" 2>/dev/null || echo 0)
if [ "$SKIPPED" -gt 0 ]; then
    echo "[PASS] second import skipped ${SKIPPED} existing objects"
    PASS=$((PASS+1))
else
    echo "[FAIL] second import expected skipped > 0, got ${SKIPPED}"
    FAIL=$((FAIL+1))
fi

# Third import with overwrite=true → should update.
assert_status "/api/v1/import" POST 200 \
    "{\"content\":$(python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))' < "$TMPFILE"),\"config_id\":0,\"overwrite\":true}"

UPDATED=$(python3 -c "import json; d=json.load(open('/tmp/nagiosql_body.json')); print(d.get('updated',0))" 2>/dev/null || echo 0)
if [ "$UPDATED" -gt 0 ]; then
    echo "[PASS] overwrite import updated ${UPDATED} objects"
    PASS=$((PASS+1))
else
    echo "[FAIL] overwrite import expected updated > 0, got ${UPDATED}"
    FAIL=$((FAIL+1))
fi

rm -f "$TMPFILE"

summary
