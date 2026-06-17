#!/usr/bin/env bash
# Timeperiod CRUD smoke tests.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

echo "=== Timeperiod Tests ==="

login admin admin123

# Create.
assert_status "/api/v1/timeperiods" POST 201 \
    '{"timeperiod_name":"smoke-24x7","alias":"24 Hours / 7 Days","ranges":[{"day":"monday","time_def":"00:00-24:00"},{"day":"sunday","time_def":"00:00-24:00"}],"active":"1","register":"1"}'
TP_ID=$(python3 -c "import json; print(json.load(open('/tmp/nagiosql_body.json')).get('id',''))")

# List.
assert_status "/api/v1/timeperiods" GET 200

# Get by ID.
assert_status "/api/v1/timeperiods/${TP_ID}" GET 200
assert_field "timeperiod_name" "smoke-24x7"

# Update.
assert_status "/api/v1/timeperiods/${TP_ID}" PUT 200 \
    '{"timeperiod_name":"smoke-24x7","alias":"All Day","ranges":[{"day":"monday","time_def":"00:00-24:00"}],"active":"1","register":"1"}'

# Delete.
assert_status "/api/v1/timeperiods/${TP_ID}" DELETE 204

# Get deleted → 404.
assert_status "/api/v1/timeperiods/${TP_ID}" GET 404

summary
