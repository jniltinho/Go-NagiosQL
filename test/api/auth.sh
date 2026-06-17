#!/usr/bin/env bash
# Auth endpoint smoke tests.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

echo "=== Auth Tests ==="

# Successful login.
login admin admin123

# /me endpoint.
assert_status "/api/v1/me" GET 200
assert_field "username" "admin"

# Bad login.
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"wrongpass"}')
if [ "$code" = "401" ]; then
    echo "[PASS] bad login → 401"
    PASS=$((PASS+1))
else
    echo "[FAIL] bad login → ${code} (expected 401)"
    FAIL=$((FAIL+1))
fi

# Token refresh using cookie.
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}/api/v1/auth/refresh" \
    -b /tmp/nagiosql_cookies.txt)
if [ "$code" = "200" ]; then
    echo "[PASS] refresh → 200"
    PASS=$((PASS+1))
else
    echo "[FAIL] refresh → ${code} (expected 200)"
    FAIL=$((FAIL+1))
fi

# Logout.
assert_status "/api/v1/auth/logout" POST 204

summary
