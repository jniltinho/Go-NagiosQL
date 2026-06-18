#!/bin/bash
# Monitora /var/lib/nagios4/rw/reload.trigger e recarrega o Nagios quando acionado.
# nagiosql API grava esse arquivo após aplicar configurações.
# Valida nagios.cfg antes de enviar SIGHUP — evita quebrar o monitoramento.

TRIGGER=/var/lib/nagios4/rw/reload.trigger
NAGIOS_CFG=/etc/nagios4/nagios.cfg
NAGIOS_BIN=/usr/sbin/nagios4

while true; do
    # -s: file exists AND is non-empty
    if [ -s "$TRIGGER" ]; then
        PID=$(pgrep -x nagios4 2>/dev/null | head -1)

        if [ -z "$PID" ]; then
            # nagios ainda não subiu — mantém o trigger para tentar no próximo ciclo
            sleep 2
            continue
        fi

        # Só limpa depois de confirmar que nagios está rodando
        > "$TRIGGER"
        echo "[reload-watcher] Reload solicitado..."

        if ! "$NAGIOS_BIN" -v "$NAGIOS_CFG" >/dev/null 2>&1; then
            echo "[reload-watcher] ERRO: $NAGIOS_CFG inválido — reload cancelado."
        else
            kill -HUP "$PID"
            echo "[reload-watcher] Config válido. SIGHUP enviado ao PID $PID"
        fi
    fi
    sleep 2
done
