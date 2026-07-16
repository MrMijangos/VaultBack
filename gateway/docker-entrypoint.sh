#!/bin/sh
set -e

# Railway asigna el puerto público de forma dinámica vía $PORT -- nginx no
# lee variables de entorno directamente, así que se sustituye aquí antes de
# arrancar. Si $PORT no está seteada (p.ej. corriendo el contenedor suelto,
# fuera de Railway), cae a 80.
: "${PORT:=80}"
export PORT

# Resolver DNS real del contenedor, para las peticiones a los upstreams de
# *.railway.internal (ver nginx.conf) -- se lee de resolv.conf en vez de
# asumir una IP fija, para no depender de la infraestructura interna
# específica de Railway.
RESOLVER=$(awk '/^nameserver/{print $2; exit}' /etc/resolv.conf)
: "${RESOLVER:=127.0.0.11}"
export RESOLVER

envsubst '${PORT} ${RESOLVER}' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf

exec nginx -g 'daemon off;'
