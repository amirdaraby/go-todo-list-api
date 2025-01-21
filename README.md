# todo list API

## How to Run

1. Generate the required keys if not already present:
   ```bash
   openssl ecparam -genkey -name prime256v1 -noout -out keys/private.pem
   openssl ec -in keys/private.pem -pubout -out keys/public.pem
   ```

2. Copy .env.example to .env
   ```bash
   cp .env.example .env
   ```

3. Run the application:
   ```bash
   docker compose up -d --build --force-recreate
   ```
