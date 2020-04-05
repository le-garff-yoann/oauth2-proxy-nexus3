# compose-gitlab-nexus3-oauth-setup

> Demonstrative setup of Sonatype Nexus 3 using GitLab as an un OAuth2 provider via [oauth2-proxy](https://github.com/oauth2-proxy/oauth2-proxy), nginx and [nexus3-proxy](nexus3-proxy/).

## TL;DR

### Prerequisites

- `docker` and `docker-compose`.
- `bash`, `curl`, `jq` and `uuidgen`.

### Setup

```bash
sudo bash setup.sh
```

- GitLab is exposed on *https://gitlab.localhost*.
- Nexus 3 is exposed on *https://nexus3.localhost* and *https://nexus3-direct.localhost*.

Passwords are stored at the end of setup in the `.secrets.env` file.

### Cleanup

```bash
sudo docker-compose down -v --remove-orphans --rmi all
```

## "Not accurate" flow 

```
********** 1↔↔ ********* 1↔↔ **************** 5↔↔ **************** 5↔↔ ***********
*        * 2↔↔ *       * 3↔↔ *              *     * nexus3-proxy *     * Nexus 3 *
* Client * 3↔↔ * Nginx * 4↔↔ * oauth2-proxy *     ****************     ***********
*        * 4↔↔ *       * 5↔↔ *              *
********** 5↔↔ *********     ****************
                       2     3
                       ↕     ↕
                       ↔↔↔↔↔ ****************
                             * GitLab (IDP) *
                             ****************
```

1. Sign in and redirect to the IDP.
2. Login and authorize the application.
3. Ask for a token.
4. Follow the callback to *oauth2-proxy* and finalize the Oauth 2.0 flow. *oauth2-proxy* then set some headers for each authorized requests.
5. *oauth2-proxy* verify authorize each request to *nexus3-proxy* . Headers provided by *oauth2-proxy* are used to keep in sync the Nexus 3 userbase.
