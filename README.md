# goredirect - HTTP Host redirection.

This aims to redirect hosts specified in `config.json` to another URL.

NOTE: Only domains will be redirected, not paths. However, you may redirect to a path.

Example `config.json`:

    {
      "example.de": "http://example.com",
      "example.fi": "http://example.com"
    }

These will redirect from example.de to http://example.com and example.fi to http://example.com

The config will load into memory and can be reloaded by running a GET query to `/reload`.

## Build

    mkdir bin
    go build -o bin/goredirect goredirect.go
