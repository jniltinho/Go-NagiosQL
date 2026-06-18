#!/bin/bash
# Monitora /var/lib/nagios4/rw/reload.trigger e recarrega o Nagios quando acionado.
# O NagiosQL grava esse arquivo ao salvar configurações (campo commandfile em tbl_configtarget).
# O dir rw/ tem grupo www-data (Debian nagios4 padrão) — Apache pode criar o arquivo.
# Valida nagios.cfg antes de enviar SIGHUP — evita quebrar o monitoramento.

TRIGGER=/var/lib/nagios4/rw/reload.trigger
NAGIOS_CFG=/etc/nagios4/nagios.cfg
NAGIOS_BIN=/usr/sbin/nagios4

while true; do
    # -s: file exists AND is non-empty (NagiosQL wrote RESTART_PROGRAM to it)
    if [ -s "$TRIGGER" ]; then
        # Truncate immediately so NagiosQL can write again without waiting
        > "$TRIGGER"
        echo "[reload-watcher] Reload solicitado pelo NagiosQL..."

        PID=$(pgrep -x nagios4 2>/dev/null | head -1)

        if [ -z "$PID" ]; then
            echo "[reload-watcher] nagios4 não está rodando"
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
