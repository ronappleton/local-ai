WIP ---
New
+83
-0

# Codex: Local AI Assistant

Codex is a minimal AI assistant that runs completely offline. It exposes a small
command line interface and an optional web UI for chatting with a language model.
The project communicates with a [llama.cpp](https://github.com/ggerganov/llama.cpp)
server and stores conversation history in a local SQLite database.

## Features

- **HTTP chat API** served via `codex serve`
- **Project management** to isolate conversations
- **Long‑term memory** stored in `memory.db`
- **Hugging Face model downloader** with `codex models`
- **Lightweight web client** under the `client/` directory
- **Tailwind CSS** is used for styling the web client
- **Simple admin area** served from the web client

## Building

Go 1.24 or newer is required. Clone the repository and run:

```bash
go build -o codex ./
```

This produces the `codex` binary in the current directory.

## Usage

Start the local HTTP API and web interface:

```bash
codex serve
```

The server listens on `http://localhost:8081` and expects a llama.cpp compatible
model server on `http://localhost:8080` for completions. A simple way to run both
services is via Docker Compose:

```bash
docker-compose up
```

The compose file includes an SMTP server under the `mail` service. The Codex
container sends email notifications through this server using the environment
variables `SMTP_ADDR` and `SMTP_FROM`.

### CLI commands

- `codex add [project] [role] [content]` – store a message in memory
- `codex serve` – launch the HTTP API and web client
- `codex models list` – browse Hugging Face models by pipeline
- `codex models download [id]` – download model files
- `codex models use [id]` – mark a downloaded model as active
- `codex models status` – show the currently active model

Run `codex [command] --help` for detailed flags.

### Admin model management

Logged in administrators can manage models from the web UI under
`/admin/models`. The interface lists Hugging Face models for each pipeline
category, supports searching by ID, and allows downloading or activating a
model. Selecting a model displays detailed metadata including tags, SHA and file
list retrieved from the Hugging Face API. A separate stats view shows the total
number of models available across Hugging Face.

## Data location

All conversation history and project metadata are kept in `memory.db` in the
working directory. Downloaded models
are stored under `models/` with state tracked in `models/state.json`.

## Running tests

The repository contains unit tests for the CLI, HTTP handlers and persistence
layer. Execute them with:

```bash
go test ./...
```

## Docker image

A multi-stage Dockerfile is provided. Build and run with:

```bash
docker build -t codex .
docker run -p 8081:8081 codex
```

When used together with the compose file, the container exposes the web UI on
`http://localhost:8081`.

## License

This project is licensed under the MIT License.
