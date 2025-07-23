#!/bin/bash

script_dir=$(dirname "${BASH_SOURCE[0]}")

nexus3_configure_rut_realm()
(
    script_name=configureRutRealm
    script_type=groovy
    script_endpoint=$1/service/rest/v1/script

    echo '{"name": "", "type": "", "content": ""}' | jq -cr \
        --arg name "$script_name" --arg type "$script_type" --arg content "$(cat "$script_dir/$script_name.$script_type")" \
        '. + {name: $name, type: $type, content: $content}' > "$script_dir/$script_name.json"

    curl -sLku "$2" -H "Content-Type: application/json" -d @"$script_dir/$script_name.json" "$script_endpoint/"

    [[ $(curl -sLku "$2" -o /dev/null -w "%{http_code}" -X POST -H "Content-Type: text/plain" \
        "$script_endpoint/$script_name/run") -eq 200 ]]
    rc=$?

    curl -sLku "$2" -X DELETE "$script_endpoint/$script_name"

    rm -Rf "$script_dir/$script_name.json"

    return $rc
)

nexus3_modify_user_password()
{
    [[ $(curl -sLku "$2:$3" -o /dev/null -w "%{http_code}" -X PUT -H "Content-Type: text/plain" -d "$4" \
        "$1/service/rest/beta/security/users/$2/change-password") -eq 204 ]]
}