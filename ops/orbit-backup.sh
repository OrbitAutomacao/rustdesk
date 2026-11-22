#!/bin/sh
set -eu
umask 077

backup_root=/opt/orbit-backups
public_key="$backup_root/public.pem"
test -s "$public_key"
mkdir -p "$backup_root/daily"
work=$(mktemp -d "$backup_root/work.XXXXXXXX")
trap 'rm -rf "$work"' EXIT
mkdir "$work/stage"

python3 - "$work/stage/rustdeskapi.db" <<'PY'
import sqlite3
import sys

source = sqlite3.connect('/opt/orbit-rustdesk-api/data/rustdeskapi.db')
target = sqlite3.connect(sys.argv[1])
source.backup(target)
assert target.execute('pragma integrity_check').fetchone()[0] == 'ok'
target.close()
source.close()
PY

cp -a /opt/orbit-rustdesk/data "$work/stage/rustdesk-data"
cp -a /opt/orbit-rustdesk/compose.yaml "$work/stage/rustdesk-compose.yaml"
cp -a /opt/orbit-rustdesk-api/compose.yaml "$work/stage/api-compose.yaml"
cp -a /opt/orbit-caddy "$work/stage/orbit-caddy"
tar -C "$work/stage" -czf "$work/archive.tar.gz" .
tar -tzf "$work/archive.tar.gz" >/dev/null

openssl rand -base64 48 > "$work/pass"
stamp=$(date -u +%Y%m%dT%H%M%SZ)
prefix="$backup_root/daily/orbit-$stamp"
openssl enc -aes-256-cbc -pbkdf2 -iter 600000 -salt \
  -in "$work/archive.tar.gz" -out "$work/archive.enc" \
  -pass "file:$work/pass"
openssl pkeyutl -encrypt -pubin -inkey "$public_key" \
  -in "$work/pass" -out "$work/pass.rsa" \
  -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha256
mv "$work/archive.enc" "$prefix.tar.gz.enc"
mv "$work/pass.rsa" "$prefix.pass.rsa"
printf 'backup_created %s\n' "$stamp"
