#!/usr/bin/env bash
# Group CRUD smoke tests (hostgroup, servicegroup, contactgroup).
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

echo "=== Group Tests ==="

login admin admin123

# -- Hostgroup --
assert_status "/api/v1/hostgroups" POST 201 \
    '{"hostgroup_name":"smoke-linux-servers","alias":"Linux Servers","active":"1","register":"1"}'
HG_ID=$(python3 -c "import json; print(json.load(open('/tmp/nagiosql_body.json')).get('id',''))")

assert_status "/api/v1/hostgroups" GET 200
assert_status "/api/v1/hostgroups/${HG_ID}" GET 200
assert_field "hostgroup_name" "smoke-linux-servers"

assert_status "/api/v1/hostgroups/${HG_ID}" PUT 200 \
    '{"hostgroup_name":"smoke-linux-servers","alias":"All Linux Servers","active":"1","register":"1"}'

assert_status "/api/v1/hostgroups/${HG_ID}" DELETE 204
assert_status "/api/v1/hostgroups/${HG_ID}" GET 404

# -- Servicegroup --
assert_status "/api/v1/servicegroups" POST 201 \
    '{"servicegroup_name":"smoke-http-checks","alias":"HTTP Checks","active":"1","register":"1"}'
SG_ID=$(python3 -c "import json; print(json.load(open('/tmp/nagiosql_body.json')).get('id',''))")

assert_status "/api/v1/servicegroups" GET 200
assert_status "/api/v1/servicegroups/${SG_ID}" DELETE 204

# -- Contactgroup --
assert_status "/api/v1/contactgroups" POST 201 \
    '{"contactgroup_name":"smoke-admins","alias":"Admin Team","active":"1","register":"1"}'
CG_ID=$(python3 -c "import json; print(json.load(open('/tmp/nagiosql_body.json')).get('id',''))")

assert_status "/api/v1/contactgroups" GET 200
assert_status "/api/v1/contactgroups/${CG_ID}" DELETE 204

summary
