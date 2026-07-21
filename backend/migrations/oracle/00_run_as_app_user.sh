#!/bin/bash
set -euo pipefail

if command -v sqlplus >/dev/null 2>&1; then
  SQLPLUS_BIN="$(command -v sqlplus)"
else
  SQLPLUS_BIN="$(ls /opt/oracle/product/*/bin/sqlplus | head -n1)"
fi

echo "Applying MetaRTLS schema as ${APP_USER}..."
"${SQLPLUS_BIN}" -s "${APP_USER}/${APP_USER_PASSWORD}@//localhost/FREEPDB1" <<SQL
WHENEVER SQLERROR EXIT SQL.SQLCODE
@/container-entrypoint-initdb.d/sql/01_schema.sql
@/container-entrypoint-initdb.d/sql/02_seed.sql
EXIT
SQL

echo "MetaRTLS schema applied."
