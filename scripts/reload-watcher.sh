#!/usr/bin/env bash
# reload-watcher.sh — polls reload.trigger and issues a Nagios graceful reload.
#
# Usage:
#   reload-watcher.sh [trigger_file] [nagios_pid_file] [interval_seconds]
#
# Defaults:
#   trigger_file   = /usr/local/nagios/var/reload.trigger
#   nagios_pid_file = /usr/local/nagios/var/nagios.lock
#   interval        = 5  (seconds between polls)
#
# The trigger file is written by go-nagiosql (TriggerReload) as a Unix timestamp.
# This script compares the modification time of the file; when it changes, a
# SIGHUP is sent to the Nagios process.
#
# IMPORTANT: Do NOT send signals via the nagios.cmd FIFO from Go — opening the
# FIFO blocks until Nagios reads it, causing a deadlock in the Go server.
# This watcher is the safe alternative.
set -euo pipefail

TRIGGER_FILE="${1:-/usr/local/nagios/var/reload.trigger}"
PID_FILE="${2:-/usr/local/nagios/var/nagios.lock}"
INTERVAL="${3:-5}"

last_mtime=""

log() { echo "[$(date '+%Y-%m-%d %H:%M:%S')] reload-watcher: $*"; }

log "started (trigger=${TRIGGER_FILE} pid_file=${PID_FILE} interval=${INTERVAL}s)"

while true; do
    if [ ! -f "$TRIGGER_FILE" ]; then
        sleep "$INTERVAL"
        continue
    fi

    current_mtime=$(stat -c '%Y' "$TRIGGER_FILE" 2>/dev/null || echo "")

    if [ -n "$current_mtime" ] && [ "$current_mtime" != "$last_mtime" ]; then
        if [ -n "$last_mtime" ]; then
            # mtime changed — send SIGHUP to Nagios.
            if [ ! -f "$PID_FILE" ]; then
                log "WARNING: PID file not found at ${PID_FILE}; skipping reload"
            else
                nagios_pid=$(cat "$PID_FILE")
                if kill -0 "$nagios_pid" 2>/dev/null; then
                    log "sending SIGHUP to Nagios (pid=${nagios_pid})"
                    kill -HUP "$nagios_pid"
                else
                    log "WARNING: Nagios pid ${nagios_pid} not running"
                fi
            fi
        fi
        last_mtime="$current_mtime"
    fi

    sleep "$INTERVAL"
done
