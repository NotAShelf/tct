# tct

**t**cp **c**onnection **t**imer (tct) is a miniscule utility to help determinee
the "optimal" number of parallel TCP requests for your network connection, for a
given target.

It performs a series of tests (following your desired configuration) by
incrementally increasing the number of parallel requests and measuring the time
taken for each test.

The optimal number is identified as the point where adding more parallel
requests does not significantly reduce the overall time taken.

## Usage

```bash
Usage:
  tct [flags]

Flags:
  -d, --delay duration   Delay between requests
  -h, --help             help for tct
  -m, --max int          Maximum number of parallel requests (default 100)
  -u, --url string       URL to fetch (default "http://example.com")
  -v, --version          version for tct
```

For example:

```bash
tct --max 200 --delay 100ms --url "http://yourtargeturl.com"
```

Replace `"http://yourtargeturl.com"` with the actual URL you wish to test
against. You can also omit the URL and use an IP address instead. For example,
`8.8.8.8` for Google or `1.1.1.1` for Cloudflare. You may notice differences
between target URLs as a result of different distances to different hosts, or
different network setups.

The `--max` parameter specifies the maximum number of parallel requests to test,
and `--delay` sets the interval between each request.

> [!TIP]
> You are strongly advised to use the delay option. I have observed high latency
> while running with the default 0 second delay, which is likely some form of
> throttling by the host. If you test against an URL that you _know_ does not
> throttle connections, then you may consider omitting `-delay`.

### Flags

- `--url`: The URL to fetch.
- `--max`: Maximum number of parallel requests to test. Default is `100`
- `--delay`: Delay between requests. Can be specified as a duration (e.g.,
  `500ms`). Default is `0` (i.e. no delay)

## Motivation

The [Nix Package Manager](https://github.com/NixOS/nix) has an option called
`http-connections` that, as per the wiki, sets the _"maximum number of parallel
TCP connections used to fetch files from binary caches and by other downloads."_

I have found this option a little obscure, however, as the default value (`25`)
has no solid motivation behind it. Is it optimal? I don't know. Can it be
optimized? Probably. tcp has been designed to help you decide a loose range to
set `http-connections` in.

Do keep in mind that this is not 100% accurate. There are many factors that may
affect the results of network related tests.

## Notes

- Remember to adjust the `-max` and `-delay` parameters based on your network
  conditions and the capabilities of the target server. tct tries not to make
  any assumptions, but the default values will not be suitable for all testing
  conditions.
- tct will try to always exit gracefully when you, e.g., kill the program with
  <kbd>Ctrl+C</kbd>.

## Contributing

Contributions are always welcome. If you have any suggestions, create an issue.
If you would like to fix my code (can't blame your for it) then create a Pull
Request.

## Reference

[Implementing graceful shutdown in go]: https://www.rudderstack.com/blog/implementing-graceful-shutdown-in-go/

- [Implementing graceful shutdown in go]
