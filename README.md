# goredirect - HTTP Host redirection.

This aims to redirect hosts specified in `config.json` to another URL.

NOTE: Only domains will be redirected, not paths. However, you may redirect to a path.

Example `config.json`:

    {
      "example.de": "http://example.com",
      "example.fi": "http://example.com"
    }

These will redirect from example.de to http://example.com and example.fi to http://example.com

The config will load into memory and can be reloaded by running a GET query to `/r/config/reload`.

In a Vagrant box it can do about 3500qps and allocates around 5MB heap.

You can set the port by running it with env PORT, example:

    PORT=8080 ./bin/goredirect

If you for some reason would like to disable the output log, set the `verbose` flag to false:

    ./bin/goredirect -verbose=false

## Build

    make get-deps
    make
