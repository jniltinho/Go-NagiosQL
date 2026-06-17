#!/usr/bin/env bash
# Contact CRUD smoke tests.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

echo "=== Contact Tests ==="

login admin admin123

# Create.
assert_status "/api/v1/contacts" POST 201 \
    '{"contact_name":"smoke-ops","alias":"Ops Team","email":"ops@example.com","host_notifications_enabled":"1","service_notifications_enabled":"1","host_notification_options":"d,u,r","service_notification_options":"w,u,c,r","active":"1","register":"1"}'
CT_ID=$(python3 -c "import json; print(json.load(open('/tmp/nagiosql_body.json')).get('id',''))")

# List.
assert_status "/api/v1/contacts" GET 200

# Get by ID.
assert_status "/api/v1/contacts/${CT_ID}" GET 200
assert_field "contact_name" "smoke-ops"

# Update.
assert_status "/api/v1/contacts/${CT_ID}" PUT 200 \
    '{"contact_name":"smoke-ops","alias":"Operations","email":"ops@example.com","host_notifications_enabled":"1","service_notifications_enabled":"1","host_notification_options":"d,u,r","service_notification_options":"w,u,c,r","active":"1","register":"1"}'

# Delete.
assert_status "/api/v1/contacts/${CT_ID}" DELETE 204

# Get deleted → 404.
assert_status "/api/v1/contacts/${CT_ID}" GET 404

summary
