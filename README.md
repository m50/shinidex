# shinidex

Requires the following envvars to be set either with a `.env` file at execution root, or in the environment.

```env
DB_PATH="file:shinidex.db" #Example for if using a local DB
TURSO_URL= #Only necessary if using Turso and not a local DB
TURSO_AUTH_TOKEN= #Only necessary if using Turso and not a local DB
AUTH_KEY= #Generate a reasonable auth key
```

## Running with Docker Compose

You can run shinidex using Docker Compose with the pre-built container from GitHub Container Registry:

Create a `docker-compose.yml` file:

```yaml
services:
  shinidex:
    image: ghcr.io/m50/shinidex:latest
    ports:
      - "1323:1323"
    environment:
      - DB_PATH=file:/app/data/shinidex.db
      - AUTH_KEY=your-secure-auth-key-here  # Generate with: openssl rand -base64 32
      # Optional Turso configuration:
      # - TURSO_URL=your-turso-url
      # - TURSO_AUTH_TOKEN=your-turso-token
    volumes:
      - ./data/shinidex.db:/app/data/shinidex.db  # Mount for persistent database storage
      - ./data/imgs/:/app/assets/imgs/
    restart: unless-stopped
```

Run the application:

```bash
# Generate a secure AUTH_KEY
echo "AUTH_KEY: $(openssl rand -base64 32)"

# Update the docker-compose.yml with your generated key, then start
docker-compose up -d
```

The application will be available at http://localhost:1323

**Note:** Make sure to replace `your-secure-auth-key-here` with a securely generated key using the command provided above.
