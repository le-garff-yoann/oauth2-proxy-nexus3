# compose-oauth2-proxy-nexus3

## TL;DR

### Prerequisites

- `docker` and `docker compose`.
- `bash`, `curl`, `jq` and `uuidgen`.

### Setup

```bash
docker compose down -v --remove-orphans

bash setup.bash
```

- GitLab is exposed on *https://gitlab.localhost*.
- Nexus 3 is exposed on *https://nexus3.localhost* and *https://nexus3-direct.localhost*.

Passwords are stored at the end of setup in the `.secrets.env` file.

### Cleanup

```bash
docker compose down -v --remove-orphans --rmi all
```
