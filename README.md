# Freezer Contents

> :warning: This isn't really intended for general use. I wrote it hastily for myself, and then tidied it up using Claude code (partly as a exercise to familiarise myself with the state of coding agents). Feel free to use it, but I make no assurances about its quality.

A webapp for keeping track of what food you have stored where (in freezers, and/or fridges, or cupboards...).

It is opinionated, and asusmes that you can have a number of containers permanently labelled with identifiers, rather than their current contents, and use this app as an index to keep track of what each container's contents and location.

This has the advantage that you only need to write labels once, and as you need to write less (just a number, rather than a date and description) you can write larger, which is helpful if you have bad handwriting.
It also forces you to record where each container is stored, which you might not bother to do if you had labelled containers with their current contents. 

You could do the same thing with a fairly simple spreasheet, but the experience of editing from a mobile phone would probably be significantly worse.


## Project Structure

The project has a backend written in [Go](https://go.dev/), and a frontend written in [Svelte](https://svelte.dev/).

The backend is compiled to a single binary, with the web app assets embedded.

It stores data in a JSON file (`contents.json`).


## Building

Use the [just](https://github.com/casey/just) command runner to build the frontend and backend.

```bash
just build
```

This will:
1. Install npm dependencies and build the Svelte UI
2. Compile the Go backend

## Running

Start the server:

```bash
./freezer-contents
``

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
