#!/usr/bin/env bash
# User management smoke tests (admin-only).
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

echo "=== User Tests ==="

login admin admin123

# Create user.
assert_status "/api/v1/users" POST 201 \
    '{"username":"smoke-user","name":"Smoke User","email":"smoke@example.com","password":"SmokePass1!","admin":"0","active":"1"}'
U_ID=$(python3 -c "import json; print(json.load(open('/tmp/nagiosql_body.json')).get('id',''))")

# List users (admin only).
assert_status "/api/v1/users" GET 200

# Get by ID.
assert_status "/api/v1/users/${U_ID}" GET 200
assert_field "username" "smoke-user"

# Change password.
assert_status "/api/v1/users/${U_ID}/password" PUT 204 \
    '{"password":"NewSmoke2!"}'

# Login with new password should work.
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"smoke-user","password":"NewSmoke2!"}')
if [ "$code" = "200" ]; then
    echo "[PASS] login with new password → 200"
    PASS=$((PASS+1))
else
    echo "[FAIL] login with new password → ${code}"
    FAIL=$((FAIL+1))
fi

# Delete user.
assert_status "/api/v1/users/${U_ID}" DELETE 204

# Cannot delete self (admin).
ADMIN_ID=$(curl -sf -H "Authorization: Bearer ${TOKEN}" "${BASE_URL}/api/v1/me" | python3 -c "import json,sys; d=json.load(sys.stdin); print(d.get('id',1))" 2>/dev/null || echo 1)
code=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE \
    -H "Authorization: Bearer ${TOKEN}" "${BASE_URL}/api/v1/users/${ADMIN_ID}")
if [ "$code" = "409" ]; then
    echo "[PASS] self-delete returns 409"
    PASS=$((PASS+1))
else
    echo "[FAIL] self-delete expected 409, got ${code}"
    FAIL=$((FAIL+1))
fi

summary
