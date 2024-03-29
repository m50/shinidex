# shinidex

Requires the following envvars to be set either with a `.env` file at execution root, or in the environment.

```env
DB_PATH="file:shinidex.db" #Example for if using a local DB
TURSO_URL= #Only necessary if using Turso and not a local DB
TURSO_AUTH_TOKEN= #Only necessary if using Turso and not a local DB
AUTH_KEY= #Generate a reasonable auth key
```