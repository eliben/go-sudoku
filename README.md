# go-sudoku

This is a toolkit for solving and generating [Sudoku
puzzles](https://en.wikipedia.org/wiki/Sudoku) in Go.

## What's inside

All the code is in a single top-level package: `sudoku`. It's broadly separated
into three parts:

* `sudoku.go`: board representation and functions for parsing boards from
  strings, emitting boards back to output and solving Sudoku puzzles. The
  basic solver uses constraint propagation and recursive search and is based on
  [Peter Norvig's old post](https://norvig.com/sudoku.html), although the Go
  code is about 100x faster than Norvig's Python (faster compiled language but
  also an optimized board representation).

  Contains additional functionality like finding _all_ the solutions of a given
  puzzle and not just a single solution.

* `generator.go`: generate valid Sudoku puzzles that have a single solution.
  The algorithm is based on a mish-mash of information found online and tweaked
  by me. Contains additional functionality like generating _symmetrical_
  Sudoku boards.

* `difficulty.go`: code to evaluate the difficulty of a given Sudoku puzzle;
  the approach was partially inspired by the paper "Sudoku Puzzles Generating:
  from Easy to Evil" by Xiang-Sun ZHANG's research group.

The `cmd` directory has two command-line tools: `generator` and `solver` that
demonstrate the use of the package.

## Testing

Some tests take a while to run, so they are excluded if the `-short` testing
flag is provided:

    $ go test -v -short ./...

## Generating printable boards

`go-sudoku` includes some rudimentary functionality to emit a Sudoku board into
a printable SVG format, like this:

![SVG board sample](/doc/genb.svg)

You can invoke the `cmd/generator` command with the `-svgout` flag to see this
in action, or use the web interface.

## Web interface

This repository includes a web interface for generating Sudoku puzzles, by
compiling the Go code to WebAssembly. To run it locally:


		$ cd cmd/wasm
		$ make build
		$ make serve

This will run a local webserver; open http://localhost:8899 in your browser
to generate puzzles!

The repository also has a GitHub actions setup to automatically deploy the web
interface to https://eliben.github.io/go-sudoku on each commit.
