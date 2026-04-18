# birdhouse

My birdhouse camera project, originally written as my capstone project at Appalachian State University.

## Running Locally

### DB

```sh
psql

\c birdhouse
```

### API

```sh
cd ./api
export CONFIG_PATH=..;
export ANTHROPIC_API_KEY=..;
export DB_PASSWORD=..;
export GOOGLE_CLIENT_ID=..;
echo '====================\nfmt' && gofmt -l -d -s -w .&& echo '====================\nrun' && go run ./...
```

### Web

```sh
cd ./web
cp .env.example .env  # set VITE_API_BASE_URL and VITE_GOOGLE_CLIENT_ID
npm install
npm run dev
```

### MQTT Broker

```sh
docker compose up -d
```
