[![tests](https://github.com/psycofdj/i-luv-grandma/actions/workflows/tests.yml/badge.svg)](https://github.com/psycofdj/i-luv-grandma/actions/workflows/tests.yml) [![linter](https://github.com/psycofdj/i-luv-grandma/actions/workflows/linter.yml/badge.svg)](https://github.com/psycofdj/i-luv-grandma/actions/workflows/linter.yml) [![coverage](https://psycofdj.github.io/i-luv-grandma/coverage-badge.svg)](https://psycofdj.github.io/i-luv-grandma/coverage.txt) [![doc](https://psycofdj.github.io/i-luv-grandma/doc/badge.svg)](https://psycofdj.github.io/i-luv-grandma/doc/index.html)


<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
- [Development](#development)
- [CI Workflows](#ci-workflows)
- [Limitations](#limitations)
<!-- markdown-toc end -->


# Introduction

I love my grand'ma, she makes the best Belgian waffles in the world.

In order to remain her favorite grand-son, I decided to help her enjoying her favorite hobby:
rotating pictures of her dearest memories.

To help her out, I wrote the `i-luv-grandma` program which takes
[pbm](https://en.wikipedia.org/wiki/Netpbm) files and rotates the pictures to a given angle.

# Installation

* from release assets
  * download assets for your architecture from [latest release](https://github.com/psycofdj/i-luv-grandma/releases)
  * extract tarball: `tar xzf i-luv-grandma_1.0.0_linux_amd64.tar.gz`
* from go install: `go install gihub.com/psycofdj/i-luv-grandma`
* from source: `CGO_ENABLED=0 go build -o i-luv-grandma -ldflags='-s -w' .`

# Usage

```
usage: i-luv-grandma [options]

Rotate pbm image by given angle. Result is written to output file.

  -angle float
        rotation of given decimal angle (positive or negative) (default 90)
  -help
        print usage
  -input string
        process given input file path, '-' for stdin (default "input.pbm")
  -output string
        write to given output file path, '-' for stdout (default "output.pbm")
  -profile string
        generate pprof profile output
  -version
        outputs version and revision informations
```

Example:

```sh
$ ./i-luv-grandma --angle 180 --input dataset/valid_j.pbm --output -

P1
6 10
000000
000000
001110
010001
010000
010000
010000
010000
010000
010000
```

- original file: ![original](./dataset/720p-orig.png?raw=true "Original file")
- with 45° rotation: ![rot45](./dataset/720p-rot45.png?raw=true "45° rotation")
- with 90° rotation: ![rot90](./dataset/720p-rot90.png?raw=true "90° rotation")
- with 180° rotation: ![rot180](./dataset/720p-rot180.png?raw=true "180° rotation")

# Development

- unit-tests
  - run: `go test ./...`

- unit-tests coverage report
  - run: `go test -cover -coverprofile cover.out -v ./... && go tool cover -func=cover.out`

- static check analysis:
  - install: https://golangci-lint.run/usage/install/#local-installation
  - run: `golangci-lint run --config .golangci.yml`

- performance analysis
  - generate profile trace: `./i-luv-grandma -profile output.pprof -input dataset/4320p.pbm -output /dev/null -angle 180`
  - inspect profile: `go tool pprof -top i-luv-grandma output.pprof`

- view documentation locally
  - install pkgsite: `go install golang.org/x/pkgsite/cmd/pkgsite@latest`
  - run pkgsite: `pkgsite`
  - open browser: `sensible-browser http://localhost:8080`

# CI Workflows

- the `release` workflow:
  - triggers on new semver tags like `v1.2.3` or `v1.2.3-rc4`
  - checks that unit-tests are passing and that code is free from linter warnings
  - run [goreleaser](https://goreleaser.com/) which:
    - generates tarballs for linux and darwin arch
    - creates new github release
    - upload tarballs to releases

- the `tests` workflow:
  - triggers on new commits
  - run unit-tests

- the `linter` workflow:
  - triggers on new commits
  - run [golangci-lint](https://github.com/golangci/golangci-lint)
  - creates github [annotations](https://github.blog/2018-12-14-introducing-check-runs-and-annotations/)
    for each issues found by the linter

- the `reports` workflow:
  - triggers on new commit for `main` branch
  - creates coverage report
    - runs unit-tests and extract coverage informations
    - creates badge file with overall total result (displayed on top of this `README.md`)
  - creates documentation report
    - generate static documentation websiteusing [godoc-static](code.rocketnine.space/tslocum/godoc-static)
    - creates badge (displayed on top of this `README.md`)
  - pushes [gh-pages](https://github.com/psycofdj/i-luv-grandma/tree/gh-pages) branch which is
    served by [Github Pages](https://pages.github.com)

# Limitations

 The current rotate implementation guaranty to preserve source image size at the cost of
 possible pixel loss for those projected outside boundaries.

 This could be a problem for my beloved grandma cause she clearly lakes basic photograph skills
 and the main subject is in the bottom right corner most of the time.

 A possible improvement would we to implement a `--resize` option that allows a different
 size in result image. It could work as follow:
 - create a bigger working space ensuring all points can be projected for any given angles
   - required space size can be computed by rotating all 4 corner pixels
 - translate source image in new space matching center of rotation
 - operate pixel rotations
