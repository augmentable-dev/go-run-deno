# go-run-deno

A small `golang` HTTP server that will run and stream the output of arbitrary `deno` scripts, supplied as JSON `POST` params.
Turn your Deno JS/TS into runnable APIs over HTTP.

## Why

Deno makes it easy to run remote scripts (any code publicly reachable over HTTP).
This proof-of-concept extends that by essentially turning `deno run https://...` into an HTTP API.

## How

The source in `main.go` should be fairly straightforward, only `stdlib` packages are used.
A `Dockerfile` is also included, you can build and run it like so:

```
docker build --tag go-run-deno:latest .
docker run -p 8000:8000 go-run-deno
```

Once the server is running, send a `POST` request with a JSON body and a `Content-Type: application/json` header.
The only required parameter in the JSON of the `POST` body is a field called `location` which is a URL pointing to JS/TS that deno will run.
It's equivalent to calling `deno run https://...`.

```
curl localhost:8000 -H "content-type: application/json" -d '{"location": "https://deno.land/std/examples/welcome.ts"}'
Welcome to Deno 🦕
```
