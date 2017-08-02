## v0.3.0 (2017-08-02)

* `ck subscribers` learned to filter subscribers via `--since`, `--until`, `--cancelled`, and `--email`.
* Also, `--reverse` will now list results in reverse order.
* Auto-generate Homebrew formula when tagging a new release.

## v0.2.0 (2017-07-30)

* Introduce `ck subscribers` and `ck version` commands.
* Add `--api-key`, `--api-secret`, and `--api-endpoint` options.
* Rename `CONVERTKIT_ENDPOINT` env var to `CONVERTKIT_API_ENDPOINT` for consistency.
* Return error if API secret is required, but not present.
* Document installation via Homebrew.

## v0.1.0 (2017-07-29)

* Initial release.
