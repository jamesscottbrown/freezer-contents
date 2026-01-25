# Freezer Contents

A webapp for inventorying food stored in one or more freezers. Tracks items across multiple containers (bowls, jars) and freezer locations.

## Building

Run the build script to compile both the frontend and backend:

```bash
./build.sh
```

This will:
1. Install npm dependencies and build the Svelte UI
2. Compile the Go backend

## Running

Start the server:

```bash
./freezer-contents
```

The server runs on port 8080 by default.

## API

### Get state

Retrieve the current inventory:

```bash
curl http://localhost:8080/state
```

### Add item

Add a new item to the inventory:

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"Name": "Soup", "Date": "2023-04-01", "Freezer": "cellar freezer", "Containers": ["bowl-2", "bowl-3"]}' \
  http://localhost:8080/add
```

### Remove item

Remove a container from the inventory:

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"Container": "bowl-1"}' \
  http://localhost:8080/remove
```

### Move item

Move a container to a different freezer:

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"Container": "1", "NewFreezer": "cellar freezer"}' \
  http://localhost:8080/move
```

## Project Structure

- `main.go`, `parse.go` - Go backend server
- `ui/` - Svelte 5 frontend application
- `contents.json` - Data storage file
