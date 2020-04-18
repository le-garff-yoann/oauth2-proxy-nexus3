# nexus3-gitlaboauth-proxy

> This service is designed to operate as a reverse-proxy between [oauth2-proxy](https://github.com/oauth2-proxy/oauth2-proxy) and Sonatype Nexus 3.

## Configuration

| ENV | Mandatory? | Default value | Description |
|-|-|-|-|
| `N3GOP_LISTEN_ON` | ☓ | 0.0.0.0:8080 | The [IP]:PORT on which the HTTP server will listen. |
| `N3GOP_SSL_INSECURE_SKIP_VERIFY` | ☓ | false | Skip SSL verifications if set to `true`. |
| `N3GOP_GITLAB_URL` | ✓ | | The GitLab URL on which OAuth operations will be performed. |
| `N3GOP_GITLAB_ACCESS_TOKEN_HEADER` | ☓ | X-Forwarded-Access-Token | The name of the HTTP header on which the GitLab OAuth *access_token* will be provided to this service. |
| `N3GOP_NEXUS3_URL` | ✓ | | The Nexus 3 URL on which sync and reverse-proxying will be performed. |
| `N3GOP_NEXUS3_ADMIN_USER` | ✓ | | A Nexus 3 **admin** user. |
| `N3GOP_NEXUS3_ADMIN_PASSWORD` | ✓ | | A Nexus 3 **admin** password. |
| `N3GOP_NEXUS3_RUT_HEADER` | ☓ | X-Forwarded-User | The name of the HTTP header used by the Rut Realm/capability (Nexus 3) for the authentication. |

### Prerequisites

#### oauth2-proxy

The `-pass-access-token` flag must be set to `true`.

#### Nexus 3

The Rut Realm/capability must be enabled and configured the use the same HTTP header as configured in via `$N3GOP_NEXUS3_RUT_HEADER`.
