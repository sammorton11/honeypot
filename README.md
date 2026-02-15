# HoneyTrap

A lightweight SSH honeypot that captures and logs unauthorized access attempts. Designed to be cheap to host, easy to deploy, and useful for security research.

## What It Does

HoneyTrap listens on port 2222 and presents a fake SSH banner to attract attackers. When someone connects, it captures their IP address, network type, and any data they send (login attempts, payloads, etc.), stores everything in a PostgreSQL database, then traps the attacker in a tarpit to waste their time.

The project has three components:

- **Honeypot** — a Go TCP server that mimics an SSH service and forwards captured attempts to the API
- **API** — a Go REST API (Chi router) that stores and serves attempt data, protected by API key authentication
- **Frontend** — a simple web dashboard served by Nginx for viewing captured attempts

Everything runs in Docker containers on an internal network, with only the honeypot (port 2222) and Nginx (port 80) exposed to the outside.

## Architecture

```
Attacker → :2222 (honeypot) → internal network → :8080 (api) → PostgreSQL
Browser  → :80   (nginx)   → internal network → :8080 (api) → PostgreSQL
```

The honeypot and API communicate over a Docker internal network (`honeynet`), so the API and database are never directly exposed to the internet.

## Quick Start

### Prerequisites

- Docker and Docker Compose
- A server or VPS (DigitalOcean droplet, Linode, etc.)

### Setup

1. Clone the repository:

```bash
git clone https://github.com/yourusername/honeytrap.git
cd honeytrap
```

2. Create a `.env` file from the example:

```bash
cp .env.example .env
```

3. Generate a secure API key and update your `.env`:

```bash
openssl rand -hex 32
```

4. Start everything:

```bash
docker compose up -d
```

5. Verify it's running:

```bash
# Check the health endpoint
curl http://localhost/api/health

# Simulate an attacker
nc localhost 2222
```

## Configuration

All configuration is done through environment variables in your `.env` file:

```env
# Database
POSTGRES_USER=honeypot
POSTGRES_PASSWORD=your_secure_password_here
POSTGRES_DB=honeypot
DATABASE_URL=postgres://honeypot:your_secure_password_here@db:5432/honeypot

# API Authentication
API_SECRET_KEY=your_generated_key_here
```

See `.env.example` for a full template.

## API Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| POST | `/api/attempt` | No | Submit a captured attempt (used by honeypot) |
| GET | `/api/attempts` | Yes | List all captured attempts |
| GET | `/api/attempt/{id}` | Yes | Get a specific attempt |
| DELETE | `/api/attempts` | Yes | Delete all attempts |
| DELETE | `/api/attempt/{id}` | Yes | Delete a specific attempt |
| GET | `/api/health` | No | Health check |

Authenticated endpoints require an `X-API-KEY` header:

```bash
curl -H "X-API-KEY: your_key_here" http://localhost/api/attempts
```

## Deployment

### DigitalOcean (Recommended)

A $6/month droplet (1 vCPU, 1GB RAM) is more than enough.

1. Create a droplet running Ubuntu 24.04
2. Install Docker: `curl -fsSL https://get.docker.com | sh`
3. Clone the repo and configure your `.env`
4. Run `docker compose up -d`
5. (Optional) Point a domain at your droplet and add Let's Encrypt

### Security Considerations for Production

- Change the default honeypot port from 2222 to 22 if you want to catch more traffic (move your real SSH to another port first)
- Set up UFW to only allow ports 22 (your real SSH), 2222 (honeypot), and 80/443 (dashboard)
- Use strong, unique values for all passwords and API keys
- Consider setting up log rotation for Docker containers
- Back up your database periodically if you want to preserve data

## Legal Disclaimer

Running a honeypot may have legal implications depending on your jurisdiction. In most places, operating a honeypot on infrastructure you own is legal, but laws vary. You are solely responsible for ensuring compliance with all applicable local, state, and federal laws.

This software is provided as-is for security research and educational purposes.

## Tech Stack

- **Go** — honeypot TCP server and REST API
- **Chi** — HTTP router with middleware
- **PostgreSQL** — attempt storage
- **Nginx** — reverse proxy and static file serving
- **Docker Compose** — orchestration

## Roadmap

- [ ] GeoIP lookup for attacker locations (MaxMind GeoLite2)
- [ ] Dashboard with attempt visualizations
- [ ] Configurable SSH banner
- [ ] Webhook notifications (Discord/Slack) for new attempts
- [ ] Export attempts to CSV/JSON

## Contributing

Contributions are welcome. Please open an issue first to discuss what you'd like to change.

## License

MIT License. See [LICENSE](LICENSE) for details.
