#!/usr/bin/env bash
# Service CRUD smoke tests.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

echo "=== Service Tests ==="

login admin admin123

# Create a host first so we can link the service.
assert_status "/api/v1/hosts" POST 201 \
    '{"host_name":"svc-test-host","alias":"Svc Test","address":"10.99.98.1","check_command":"check-host-alive","active":"1","register":"1"}'
HOST_ID=$(python3 -c "import json; print(json.load(open('/tmp/nagiosql_body.json')).get('id',''))")

# Create service.
assert_status "/api/v1/services" POST 201 \
    '{"service_description":"PING","config_name":"svc-test-host","check_command":"check_ping!100,20%!500,60%","active":"1","register":"1"}'
SVC_ID=$(python3 -c "import json; print(json.load(open('/tmp/nagiosql_body.json')).get('id',''))")

# List services — filter by config_name.
assert_status "/api/v1/services?config_name=svc-test-host" GET 200

# Get by ID.
assert_status "/api/v1/services/${SVC_ID}" GET 200
assert_field "service_description" "PING"

# Delete service then host.
assert_status "/api/v1/services/${SVC_ID}" DELETE 204
assert_status "/api/v1/hosts/${HOST_ID}" DELETE 204

summary
