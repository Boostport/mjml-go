# mjml-go
[![Go Reference](https://pkg.go.dev/badge/github.com/Boostport/mjml-go.svg)](https://pkg.go.dev/github.com/Boostport/mjml-go)
[![Tests Status](https://github.com/Boostport/mjml-go/workflows/Tests/badge.svg)](https://github.com/Boostport/mjml-go)
[![Test Coverage](https://api.codeclimate.com/v1/badges/cbb1efa66cb148be5cb8/test_coverage)](https://codeclimate.com/github/Boostport/mjml-go/test_coverage)

Compile [MJML](https://mjml.io/) into HTML directly in your Go application!

## Why?
[MJML](https://github.com/mjmlio/mjml) is a JavaScript library. In order to use it with other languages,
the usual approach is to wrap the library in a Node.js HTTP server and provide an endpoint through which
applications not written in JavaScript can make HTTP requests to compile MJML into HTML.

This approach poses some challenges, for example, if MJML is upgraded to a new major version in
the deployed Node.js servers, applications calling these servers will need to be upgraded in a synchronized
manner to avoid incompatibilities. In addition, running these extra servers introduces extra moving parts
and the network into the mix.

This is why we built `mjml-go` and created an idiomatic Go API to compile MJML into HTML directly in Go applications that
can be deployed as a single Go binary.

## How?
We wrote a [simple JavaScript wrapper](js/src) that wraps around the MJML library by accepting input and returning output
using JSON. This wrapper is then bundled using webpack and compiled into a WebAssembly module using Suborbital's [Javy fork](https://github.com/suborbital/javy),
a Javascript to WebAssembly compiler. The WebAssembly module is then compressed using Brotli to yield a 10x reduction in 
file size.

During runtime, the module is decompressed and loaded into a [Wazero](https://github.com/tetratelabs/wazero) runtime 
on application start up to accept input in order to compile MJML into HTML.

### Workers
As WebAssembly modules compiled using Javy are not thread-safe and cannot be called concurrently, the library maintains
a pool of 1 to 10 instances to perform compilations. Idle instances are automatically destroyed and will be re-created when
they are needed. This means that the library is thread-safe and you can use it concurrently in multiple goroutines.

## Example
```go
func main() {
	
	input := `<mjml><mj-body><mj-section><mj-column><mj-divider border-color="#F45E43"></mj-divider><mj-text font-size="20px" color="#F45E43" font-family="helvetica">Hello World</mj-text></mj-column></mj-section></mj-body></mjml>`
	
	output, err := mjml.ToHTML(context.Background(), input, mjml.WithMinify(true))
	
	var mjmlError mjml.Error
	
	if errors.As(err, &mjmlError){
	    fmt.Println(mjmlError.Message)
	    fmt.Println(mjmlError.Details)	
	}
	
	fmt.Println(output)
}
```

## Options
The library provides a complete list of options to customize the MJML compilation process including options for
`html-minifier`, `js-beautify` and `juice`.

These are all exposed via an idiomatic Go API and a complete list can be found in the [Go documentation](https://pkg.go.dev/github.com/Boostport/mjml-go).

### Defaults
If beautify and minify are enabled, but no options were passed in, the library defaults to using the same defaults
as the MJML CLI application:

For minify:

| option                  | value   |
|-------------------------|---------|
| `CaseSensitive`         | `true`  |
| `CollapseWhitespace`    | `true`  |
| `MinifyCSS`             | `false` |
| `RemoveEmptyAttributes` | `true`  |

For beautify:

| option                     | value   |
|----------------------------|---------|
| `EndWithNewline`           | `true`  |
| `IndentSize`               | `2`     |
| `PreserveNewlines`         | `false` |
| `WrapAttributesIndentSize` | `2`     |

## Limitations
The WebAssembly module is not able to access the filesystem, so `<mj-include>` tags are ignored. The solution is to
flatten your templates during development and pass the flattened templates to `mjml.ToHTML()`.

This [example](https://github.com/mjmlio/mjml/issues/2465#issuecomment-1109515536) provides a good starting point to
create a Node.js script to do this:
```javascript
import mjml2html from 'mjml' // load default component
import components from 'mjml-core/lib/components.js'
import Parser from 'mjml-parser-xml'
import jsonToXML from 'mjml-core/lib/helpers/jsonToXML.js'

const xml = `<mjml>...</mjml>`

const mjml = Parser(xml, {
      components,
      filePath: '.',
      actualPath: '.'
    })

console.log(JSON.stringify(mjml))
console.log(jsonToXML(mjml))
```

## Differences from the MJML JavaScript library
- Beautify and minify will be removed from the library in [MJML5](https://github.com/mjmlio/mjml/pull/2204) and will be
moved into the MJML CLI. Therefore, to prepare for this move, the [wrapper](js/src) imports `html-minifier`
and `js-beautify` directly to support minifying and beautifying the output.
- In the current implementation of mjml, it is not possible to customize the output of `js-beautify`. In this library,
we have exposed those options.

## Benchmarks
We are benchmarking against a very [minimal Node.js server](js/src/server.js) serving a single API endpoint.
```
goos: linux
goarch: amd64
pkg: github.com/Boostport/mjml-go
cpu: 12th Gen Intel(R) Core(TM) i7-12700F
BenchmarkNodeJS/black-friday-20                      594           1865068 ns/op
BenchmarkNodeJS/one-page-20                          288           3978085 ns/op
BenchmarkNodeJS/reactivation-email-20                198           5969088 ns/op
BenchmarkNodeJS/real-estate-20                       153           7644823 ns/op
BenchmarkNodeJS/recast-20                            180           6747342 ns/op
BenchmarkNodeJS/receipt-email-20                     344           3396417 ns/op
BenchmarkMJMLGo/black-friday-20                       19          59296864 ns/op
BenchmarkMJMLGo/one-page-20                            8         136529250 ns/op
BenchmarkMJMLGo/reactivation-email-20                  9         121438942 ns/op
BenchmarkMJMLGo/real-estate-20                         4         265611426 ns/op
BenchmarkMJMLGo/recast-20                              5         213318243 ns/op
BenchmarkMJMLGo/receipt-email-20                       9         117782824 ns/op
PASS
ok      github.com/Boostport/mjml-go    31.910s
```

In its current state the Node.js implementation is significantly faster than `mjml-go`. However, with improvements to
Wazero (in particular [tetratelabs/wazero#618](https://github.com/tetratelabs/wazero/issues/618) and [tetratelabs/wazero#179](https://github.com/tetratelabs/wazero/issues/179)),
module instantiation times should see great improvement, reducing worker spin-up times and improving the compilation performance.

Also, we should see improvements from Javy improve these numbers as well.

## Development

### Run tests
You can run tests using docker by running `docker compose run test` from the root of the repository.

### Run benchmarks
From the root of the repository, run `go test -bench=. ./...`. Alternatively, you can run them in a docker container:
`docker compose run benchmark`

### Compile WebAssembly module and build Node.js test server
Run `docker compose run build-js` from the root of the repository.

## Other languages
Since the MJML library is compiled into a WebAssembly module, it should be relatively easy to take the compiled module and
drop it into languages with WebAssembly environments.

If you've created a library for another language, please let us know, so that we can add it to this list!