#!/bin/bash
# Monitora /usr/local/nagios/var/reload.trigger e recarrega o Nagios quando acionado.
# O NagiosQL grava esse arquivo ao salvar configurações (via commandfile no tbl_configtarget).
# Valida nagios.cfg antes de enviar SIGHUP — evita quebrar monitoramento por config inválida.
# Ref: NAGIOSQL_DEBIAN_PACKAGING.md §Validation Workflow

TRIGGER=/usr/local/nagios/var/reload.trigger
NAGIOS_CFG=/usr/local/nagios/etc/nagios.cfg
NAGIOS_BIN=/usr/local/nagios/bin/nagios

while true; do
    if [ -f "$TRIGGER" ]; then
        rm -f "$TRIGGER"
        echo "[reload-watcher] Reload solicitado pelo NagiosQL..."

        # pgrep em vez de lock file: mais robusto — independente do lock_file path
        # compilado (Nagios 4.x default é /run/nagios.lock, dir root-owned)
        PID=$(pgrep -x nagios 2>/dev/null | head -1)

        if [ -z "$PID" ]; then
            echo "[reload-watcher] Nagios não está rodando"
        elif ! "$NAGIOS_BIN" -v "$NAGIOS_CFG" >/dev/null 2>&1; then
            echo "[reload-watcher] ERRO: $NAGIOS_CFG inválido — reload cancelado." \
                 "Corrija o erro no NagiosQL antes de tentar novamente."
        else
            kill -HUP "$PID"
            echo "[reload-watcher] Config válido. SIGHUP enviado ao PID $PID"
        fi
    fi
    sleep 2
done
