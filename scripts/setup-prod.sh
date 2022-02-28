APPS_FOLDER=/var/lib
APP_NAME=sport
APP_VERSION=v0.1.0;
DOMAIN=${DOMAIN:?"DOMAIN is required as env variable (e.g. mondomain.com)"}

apt install awscli sqlite3

mkdir -p ${APPS_FOLDER}/${APP_NAME}/${APP_VERSION} ${APPS_FOLDER}/${APP_NAME}/tmp/sessions ${APPS_FOLDER}/${APP_NAME}/tmp/uploads;

chown -R app:app ${APPS_FOLDER}/${APP_NAME}

cat <<EOF | curl -X POST -H 'Content-Type: application/json' -d@- localhost:2019/config/apps/http/servers/srv0/routes/
{
  "@id": "sport",
  "match": [{ "host": ["${DOMAIN}"] }],
  "handle": [
    {
      "handler": "headers",
      "response": {
        "set": {
          "Referrer-Policy": ["same-origin"],
          "X-Content-Type-Options": ["nosniff"],
          "X-Frame-Options": ["DENY"],
          "X-Xss-Protection": ["1; mode=block"]
        }
      }
    },
    {
      "handler": "reverse_proxy",
      "upstreams": [{"dial": "localhost:10001"}]
    }
  ]
}
EOF

echo "Copy/Paste env variables. Press <C-d> when done"
cat > ${APPS_FOLDER}/${APP_NAME}/.env;

cat <<CODE > /usr/bin/sport
#!/usr/bin/env bash

set -o errexit
set -o pipefail

source ${APPS_FOLDER}/${APP_NAME}/.env

case "\$1" in
  start)
    ${APPS_FOLDER}/${APP_NAME}/sport
    ;;
  dump)
    cat <<EOF | sqlite3 ${APPS_FOLDER}/${APP_NAME}/sport.sqlite | bzip2 --best | env AWS_ACCESS_KEY_ID=\${TRAX_BACKUP_AWS_ACCESS_KEY_ID} AWS_SECRET_ACCESS_KEY=\${TRAX_BACKUP_AWS_SECRET_ACCESS_KEY} AWS_REGION=\${TRAX_BACKUP_AWS_REGION} aws s3 cp - s3://\${TRAX_BACKUP_AWS_BUCKET}/dumps/sport/\$(date '+%Y-%m-%d').sql.bzip
.output stdout
.dump
.exit
EOF
    ;;
  *)
    >&2 echo "command does not exist. possible commands: start, dump"
    exit 1;
esac
CODE

chmod u+x /usr/bin/sport
chown app:app /usr/bin/sport

cat <<EOF > /etc/systemd/system/sport.service
[Unit]
Description=Sport
Documentation=https://github.com/lonepeon/sport
After=network.target network-online.target
Requires=network-online.target

[Service]
Type=simple
User=app
Group=app
ExecStart=/usr/bin/sport start
Restart=always
TimeoutStopSec=5s
LimitNOFILE=1048576
LimitNPROC=512
PrivateTmp=true
ProtectSystem=full

[Install]
WantedBy=multi-user.target
EOF

cat <<EOF > /etc/systemd/system/sport-backup.service
[Unit]
Description=Sport Backup
Wants=sport-backup.timer
After=network.target network-online.target
Requires=network-online.target

[Service]
Type=oneshot
User=app
Group=app
ExecStart=/usr/bin/sport dump

[Install]
WantedBy=multi-user.target
EOF

cat <<EOF > /etc/systemd/system/sport-backup.timer
[Unit]
Description=Dump sport database every week
Requires=sport-backup.service

[Timer]
Unit=sport-backup.service
OnCalendar=Sun *-*-* 14:00:00

[Install]
WantedBy=timers.target
EOF


systemctl enable --now sport-backup.service
systemctl enable --now sport-backup.timer
systemctl enable --now sport.service
