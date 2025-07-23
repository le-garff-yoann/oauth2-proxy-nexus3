#!/bin/bash

set -eo pipefail

# shellcheck disable=SC1091
. gitlab/configure.bash
# shellcheck disable=SC1091
. nexus3/configure.bash

wait_for()
(
    echo "➡️  Waiting for $1 to be up."

    set +e

    i=0
    until [[ $(curl -sLo /dev/null -kw "%{http_code}" "$1") == 2* ]]
    do
        sleep 1

        if [[ $i -gt ${2:-300} ]]
        then
            return 1
        fi

        # shellcheck disable=SC2219
        let i++
    done
)

GITLAB_ROOT_PASSWORD=$(uuidgen)
OAUTH2_PROXY_CLIENT_ID=$(uuidgen)
OAUTH2_PROXY_CLIENT_SECRET=$(uuidgen)

export GITLAB_ROOT_PASSWORD OAUTH2_PROXY_CLIENT_ID OAUTH2_PROXY_CLIENT_SECRET

echo "➡️  Bring the whole stack up."

docker compose up -d --build

wait_for https://gitlab.localhost
wait_for https://nexus3-direct.localhost

echo "➡️  Configure the oauth2-proxy/nexus3 application in GitLab."

oauth2proxynexus3_application=$(gitlab_configure_oauth2proxynexus3_application "$(docker compose ps -q gitlab)")

O2PN3_NEXUS3_ADMIN_USER="admin"
O2PN3_NEXUS3_ORIGINAL_ADMIN_PASSWORD=$(docker compose exec nexus3 cat /nexus-data/admin.password)
O2PN3_NEXUS3_ADMIN_PASSWORD=$(uuidgen)

export O2PN3_NEXUS3_ADMIN_USER O2PN3_NEXUS3_ORIGINAL_ADMIN_PASSWORD \
    O2PN3_NEXUS3_ADMIN_PASSWORD

echo "➡️  Change Nexus 3's default admin password."

nexus3_modify_user_password https://nexus3-direct.localhost \
    "$O2PN3_NEXUS3_ADMIN_USER" "$O2PN3_NEXUS3_ORIGINAL_ADMIN_PASSWORD" "$O2PN3_NEXUS3_ADMIN_PASSWORD" 1>/dev/null

echo "➡️  Configure Nexus 3's Rut realm."

nexus3_configure_rut_realm https://nexus3-direct.localhost \
    "$O2PN3_NEXUS3_ADMIN_USER:$O2PN3_NEXUS3_ADMIN_PASSWORD" 1>/dev/null

OAUTH2_PROXY_CLIENT_ID=$(echo "$oauth2proxynexus3_application" | jq -cr .uid)
OAUTH2_PROXY_CLIENT_SECRET=$(echo "$oauth2proxynexus3_application" | jq -cr .secret)
OAUTH2_PROXY_COOKIE_NAME=oauth2-proxy-nexus3

export OAUTH2_PROXY_COOKIE_NAME

echo "➡️  Reconfigure the stack."

docker compose up -d --build --force-recreate oauth2-proxy oauth2-proxy-nexus3
docker compose restart

wait_for https://gitlab.localhost
wait_for https://nexus3.localhost/ping
wait_for https://nexus3-direct.localhost

cat > .secrets.env <<EOF
GITLAB_ROOT_PASSWORD=$GITLAB_ROOT_PASSWORD
NEXUS3_ADMIN_PASSWORD=$O2PN3_NEXUS3_ADMIN_PASSWORD
EOF

echo "➡️  The stack is ready to be tested."

cat .secrets.env
