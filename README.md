# Guild Manager for World of Warcraft

A web application for managing WoW guilds with features like roster management, raid scheduling, and character gear tracking. Supports self-hosting and mobile packaging.

## Features
- Battle.net OAuth2 authentication (WIP)
- Daily character data sync from Blizzard API (WIP)
- Raid calendar with attendance tracking (WIP)
- Mobile push notifications (WIP)
- Admin management system (WIP)
- Cross-platform compatibility (WIP)

## Quick Start with Docker

### Prerequisites
- Docker 20.10+
- Docker Compose 2.20+

1. Clone repository:
```bash
git clone https://github.com/GFerreiroS/guild-manager.git
cd guild-manager
```

2. Configure environment:
```bash
cp .env.example .env
# Edit .env with your Battle.net API keys
```

3. Start services:
```bash
docker-compose up -d
```

Access:
- Frontend: http://localhost:80
- Backend API: http://localhost:8080

## Manual Installation

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Nginx (or reverse proxy)

1. Database setup:
```bash
createdb guild_manager
createuser guild_admin
psql guild_manager -c "ALTER USER guild_admin WITH PASSWORD 'your_password';"
```

2. Run migrations:
```bash
cd backend
go run cmd/server/main.go migrate
```

3. Start backend:
```bash
BNET_CLIENT_ID=your_id BNET_CLIENT_SECRET=your_secret go run cmd/server/main.go
```

4. Frontend setup:
```nginx
# nginx.conf
server {
    listen 80;
    root /path/to/frontend/templates;
    
    location /api/ {
        proxy_pass http://localhost:8080;
    }
}
```

## Configuration
Required environment variables:
```ini
POSTGRES_USER=guild_admin
POSTGRES_PASSWORD=your_db_password
BNET_CLIENT_ID=your_bnet_id
BNET_CLIENT_SECRET=your_bnet_secret
SESSION_SECRET=complex_secret_here
```

## Contributing
PRs welcome! Please follow:
1. Fork repository
2. Create feature branch
3. Submit PR with description

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.