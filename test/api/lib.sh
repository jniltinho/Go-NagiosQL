#!/usr/bin/env bash
# Shared helpers for NagiosQL API smoke tests.
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8081}"
TOKEN=""
PASS=0
FAIL=0

login() {
    local user="${1:-admin}"
    local pass="${2:-admin123}"
    local body
    body=$(curl -sf -X POST "${BASE_URL}/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"${user}\",\"password\":\"${pass}\"}" \
        -c /tmp/nagiosql_cookies.txt)
    TOKEN=$(echo "$body" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    if [ -z "$TOKEN" ]; then
        echo "[FAIL] login: no access_token in response: $body" >&2
        FAIL=$((FAIL+1))
        return 1
    fi
    echo "[PASS] login as ${user}"
    PASS=$((PASS+1))
    export TOKEN
}

api_get() {
    local path="$1"
    curl -sf -H "Authorization: Bearer ${TOKEN}" "${BASE_URL}${path}"
}

api_post() {
    local path="$1"
    local data="$2"
    curl -sf -X POST -H "Authorization: Bearer ${TOKEN}" \
        -H "Content-Type: application/json" \
        -d "$data" "${BASE_URL}${path}"
}

api_put() {
    local path="$1"
    local data="$2"
    curl -sf -X PUT -H "Authorization: Bearer ${TOKEN}" \
        -H "Content-Type: application/json" \
        -d "$data" "${BASE_URL}${path}"
}

api_delete() {
    local path="$1"
    curl -sf -X DELETE -H "Authorization: Bearer ${TOKEN}" "${BASE_URL}${path}"
}

assert_status() {
    local url="$1"
    local method="${2:-GET}"
    local want_code="${3:-200}"
    local data="${4:-}"
    local code
    if [ "$method" = "GET" ]; then
        code=$(curl -s -o /tmp/nagiosql_body.json -w "%{http_code}" \
            -H "Authorization: Bearer ${TOKEN}" "${BASE_URL}${url}")
    else
        code=$(curl -s -o /tmp/nagiosql_body.json -w "%{http_code}" \
            -X "$method" -H "Authorization: Bearer ${TOKEN}" \
            -H "Content-Type: application/json" \
            -d "$data" "${BASE_URL}${url}")
    fi
    if [ "$code" = "$want_code" ]; then
        echo "[PASS] ${method} ${url} → ${code}"
        PASS=$((PASS+1))
    else
        echo "[FAIL] ${method} ${url} → ${code} (expected ${want_code})"
        echo "  body: $(cat /tmp/nagiosql_body.json)"
        FAIL=$((FAIL+1))
    fi
}

assert_field() {
    local field="$1"
    local want="$2"
    local got
    got=$(python3 -c "import json,sys; d=json.load(open('/tmp/nagiosql_body.json')); print(d.get('$field',''))" 2>/dev/null || echo "")
    if [ "$got" = "$want" ]; then
        echo "[PASS] field ${field}=${want}"
        PASS=$((PASS+1))
    else
        echo "[FAIL] field ${field}: got=${got} want=${want}"
        FAIL=$((FAIL+1))
    fi
}

summary() {
    echo ""
    echo "=== Test Summary ==="
    echo "PASS: ${PASS}  FAIL: ${FAIL}"
    if [ "$FAIL" -gt 0 ]; then
        exit 1
    fi
}
