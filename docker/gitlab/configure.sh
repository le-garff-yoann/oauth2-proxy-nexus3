#!/bin/bash

gitlab_configure_oauth2proxynexus3_application()
{
    app_name=oauth2-proxy-nexus3

    cat <<EOF | docker exec -i $1 gitlab-rails r -
Doorkeeper::Application
    .new(name: '$app_name', redirect_uri: 'https://nexus3.localhost/oauth2/callback', scopes: "openid\nprofile\nemail")
    .save

puts Doorkeeper::Application
    .find_by(name: '$app_name')
    .to_json
EOF

    return $RC
}
