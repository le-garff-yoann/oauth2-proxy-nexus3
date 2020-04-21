#!/bin/bash

. gitlab/configure.sh
. nexus3/configure.sh

die() {
    echo $@

    exit 1
}

wait_for()
{
    echo "Waiting for $1 to be up."

    i=0
    until [[ $(curl -sLo /dev/null -kw '%{http_code}' $1) == 2* ]]
    do
        sleep 1

        [[ $i -gt ${2:-300} ]] && return 1

        let i++
    done
}

export GITLAB_ROOT_PASSWORD=$(uuidgen)

docker-compose up -d --build || die 'docker-compose returned an error.'

wait_for https://gitlab.localhost && \
wait_for https://nexus3-direct.localhost

echo 'Configure the oauth2-proxy/nexus3 application in GitLab.'
oauth2proxynexus3_application=$(gitlab_configure_oauth2proxynexus3_application $(docker-compose ps -q gitlab 2>/dev/null) \
    || die 'Error while configuring the oauth2-proxy/nexus3 application in GitLab.')

export \
    O2PN3_NEXUS3_ADMIN_USER=admin \
    O2PN3_NEXUS3_ORIGINAL_ADMIN_PASSWORD=$(docker-compose exec nexus3 cat /nexus-data/admin.password) \
    O2PN3_NEXUS3_ADMIN_PASSWORD=$(uuidgen)

echo "Change Nexus 3's default admin password."
nexus3_modify_user_password https://nexus3-direct.localhost \
    $O2PN3_NEXUS3_ADMIN_USER $O2PN3_NEXUS3_ORIGINAL_ADMIN_PASSWORD $O2PN3_NEXUS3_ADMIN_PASSWORD 1>/dev/null \
    || die "Error while configuring Nexus 3's default admin password."

echo "Configure Nexus 3's Rut realm."
nexus3_configure_rut_realm https://nexus3-direct.localhost \
    "$O2PN3_NEXUS3_ADMIN_USER:$O2PN3_NEXUS3_ADMIN_PASSWORD" 1>/dev/null \
    || die "Error while configuring Nexus 3's Rut realm."

export \
    OAUTH2_PROXY_CLIENT_ID=$(echo $oauth2proxynexus3_application | jq -cr .uid) \
    OAUTH2_PROXY_CLIENT_SECRET=$(echo $oauth2proxynexus3_application | jq -cr .secret) \
    OAUTH2_PROXY_COOKIE_NAME=oauth2-proxy-nexus3

docker-compose up -d --build --force-recreate oauth2-proxy oauth2-proxy-nexus3 \
    || die 'docker-compose returned an error.'

docker-compose restart || die 'docker-compose returned an error.'

wait_for https://gitlab.localhost && \
wait_for https://nexus3.localhost/ping && \
wait_for https://nexus3-direct.localhost

cat <<EOF > .secrets.env
GITLAB_ROOT_PASSWORD=$GITLAB_ROOT_PASSWORD
NEXUS3_ADMIN_PASSWORD=$O2PN3_NEXUS3_ADMIN_PASSWORD
EOF

cat .secrets.env
