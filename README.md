# Gist Sync

A daemon that syncs your GitHub gists to a self-hosted Opengist instance.

## Features

- One-way sync from GitHub to Opengist
- Periodic synchronization (configurable interval)
- Clones gists using Git operations (preserves history)
- All synced gists are marked as private on Opengist
- Dockerized with Docker Compose support

## Prerequisites

- Docker and Docker Compose
- GitHub Personal Access Token with `gist` scope
- Running Opengist instance with an access token

## Quick Start

1. Clone this repository
2. Copy `.env.example` to `.env` and configure:
   ```bash
   cp .env.example .env
   ```
3. Edit `.env` with your credentials:
   - `GITHUB_TOKEN`: Your GitHub Personal Access Token
   - `GITHUB_USERNAME`: Your GitHub username
   - `OPENGIST_URL`: URL of your Opengist instance
   - `OPENGIST_USERNAME`: Your Opengist username
   - `OPENGIST_TOKEN`: Your Opengist access token
4. Start the services:
   ```bash
   docker-compose up -d
   ```

## Configuration

All configuration is done via environment variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `GITHUB_TOKEN` | Yes | - | GitHub Personal Access Token |
| `GITHUB_USERNAME` | Yes | - | Your GitHub username |
| `OPENGIST_URL` | Yes | - | Opengist instance URL |
| `OPENGIST_USERNAME` | Yes | - | Your Opengist username |
| `OPENGIST_TOKEN` | Yes | - | Opengist access token |
| `SYNC_INTERVAL_MINUTES` | No | 30 | Sync interval in minutes |
| `LOG_LEVEL` | No | info | Logging level (debug, info, warn, error) |

## How It Works

1. The syncer fetches all your gists from GitHub using the REST API
2. For each gist, it clones the Git repository from GitHub
3. It pushes the gist to your Opengist instance using force push
4. All gists are synced as private on Opengist
5. The process repeats every `SYNC_INTERVAL_MINUTES`

## Development

### Local Build

```bash
# Build
go build -o gist-sync ./cmd/syncd

# Run with env vars
GITHUB_TOKEN=xxx GITHUB_USERNAME=xxx ... ./gist-sync
```

### Docker Build

```bash
docker build -t gist-sync .
docker run --env-file .env gist-sync
```

## License

This project is licensed under the GNU General Public License v3.0 (GPL-3.0) - see the [LICENSE](LICENSE) file for details.
