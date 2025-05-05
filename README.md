# Biathlon app

## Set up

Clone the repository

Set up .env file by example
```bash
FILE_CONFIG=absolute/path/to/config/json
FILE_EVENTS=absolute/path/to/events/txt
```

Start Docker Engine

## Start

Change directory to directory with docker-compose file
```bash
cd biathlon-app
```

Run application
```bash
docker compose up
```

## Output

Logs of steps competitors
```bash
[10:28:29.629] The target(1) has been hit by competitor(5)
```
Final table
```bash
Resulting table:
biathlon_container  | [00:25:16.853] 2 [{00:12:38.243, 4.616} {00:12:38.610, 4.614}] {00:01:40.000, 1.500} 8/10
```