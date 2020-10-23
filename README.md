<h1 align="center">gochk</h1>

<p align="center">
  <a href="https://github.com/resotto/gochk/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-GPL%20v3.0-brightgreen.svg" /></a>
</p>

<p align="center">
  Static Dependency Analysis Tool for Go Files
</p>

<p align="center">
  <img src="https://user-images.githubusercontent.com/19743841/97001249-0983ff80-1573-11eb-818f-9bdbffe8f762.gif">
</p>

---

What is gochk?

- gochk checks whether .go files violate [Clean Architecture The Dependency Rule](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html#the-dependency-rule) or not, and prints its results.

  > This rule says that source code dependencies can only point inwards. Nothing in an inner circle can know anything at all about something in an outer circle.

Why gochk?

- ZERO Dependency
- Simple & Easy-to-Read Outputs

---

## Table of Contents

- [Getting Started](#getting-started)
- [Installation](#installation)
- [How gochk works](#how-gochk-works)
- [How to see results](#how-to-see-results)
- [Configuration](#configuration)
- [Unit Testing](#unit-testing)
- [Performance Test](#performance-test)
- [Customization](#customization)
- [Build](#build)
- [Feedback](#feedback)
- [License](#license)

## Getting Started

```zsh
go get -u github.com/resotto/gochk
cd ${GOPATH}/src/github.com/resotto/gochk
```

Please edit paths of `dependencyOrders` in `gochk/configs/config.json` according to your dependency rule, whose smaller index value means outer circle.

```json
"dependencyOrders": ["external", "adapter", "application", "domain"],
```

And then, let's gochk your target path with `-t`:

```zsh
go run cmd/gochk/main.go -t=${YourTargetPath}
```

If you have [Goilerplate](https://github.com/resotto/goilerplate), you can also gochk it:

```zsh
go run cmd/gochk/main.go -t=../goilerplate
```

If your current working directory is not in gochk root `${GOPATH}/src/github.com/resotto/gochk`, you must specify the location of the `config.json` with `-c`:

```zsh
cd internal
go run ../cmd/gochk/main.go -t=../../goilerplate -c=../configs/config.json
```

## Installation

First of all, let's check whether `GOPATH` is set or not:

```zsh
go env GOPATH
```

And then, please check whether `${GOPATH}/bin` is included in your `$PATH`:

```zsh
echo $PATH
```

Finally, please install gochk:

```zsh
cd cmd/gochk
go install
```

## How gochk works

### Prerequisites

- **Please format all .go files with one of the following format tools in advance, which means only one import statement in a .go file**.
  - goimports
  - goreturns
  - gofumports
- If you have files with following file path or import path, gochk might not work well.
  - The path including the two directory names specified in `dependencyOrders` in `gochk/configs/config.json`.
    - For example, if you have the path `app/external/adapter/service` and want to handle this path as what is in `adapter`, and `dependencyOrders = ["external", "adapter"]`, the index of the path will be `0` (not `1`).

### What gochk does

gochk checks whether .go files violate [Clean Architecture The Dependency Rule](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html#the-dependency-rule) or not, and prints its results.

> This rule says that source code dependencies can only point inwards. Nothing in an inner circle can know anything at all about something in an outer circle.

For example, if an usecase in "Use Cases" imports (depends on) what is in "Controllers/Gateways/Presenters", it violates dependency rule.

<p align="center">
  <img src="https://user-images.githubusercontent.com/19743841/93830264-afa9c480-fcaa-11ea-9589-7c5308c291f4.jpg">
</p>
<p align="center">
  <a href="https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html">The Clean Architecture</a>
</p>

### Check Logic

Firstly, gochk fetchs the file path and gets the index of `dependencyOrders` in `gochk/configs/config.json` if one of them is included in the file path.

Second, gochk reads the file, parses import paths, and also gets the index of `dependencyOrders` if matched.

And then, gochk compares those indices and detects violation **if the index of the import path is smaller than that of the file path**.

For example, if you have a file `app/application/usecase/xxx.go` with import path `"app/adapter/service"` and `dependencyOrders = ["adapter", "application"]`, the index of the file is `1` and the index of its import is `0`. Therefore, the file violates dependency rule since the inequality `0 (the index of the import path) < 1 (the index of the file path)` is established.

## How to see results

### Quick Check

You can check whether there are violations or not quickly by looking at the end of results.

If you see `Dependencies which violate dependency orders found!`, there are violations!🚨

```
2020/10/19 23:37:03 Dependencies which violate dependency orders found!
```

If you see the following AA, congrats! there are no violations🎉

```
2020/10/19 23:57:25 No violations
    ________     _______       ______    __     __    __   _ _
   /  ______\   /  ___  \     /  ____\  |  |   |  |  |  | /   /
  /  /  ____   /  /   \  \   /  /       |  |___|  |  |  |/   /
 /  /  |_   | |  |     |  | |  |        |   ___   |  |      /
 \  \    \  | |  |     |  | |  |        |  |   |  |  |  |\  \
  \  \___/  /  \  \___/  /   \  \_____  |  |   |  |  |  | \  \
   \_______/    \_______/     \_______\ |__|   |__|  |__|  \__\
```

### Result types

gochk displays each result type in a different color by default:

- ![#008080](https://via.placeholder.com/15/008080/000000?text=+) None
  - which means there are imports irrelevant to dependency rule or no imports at all
- ![#008000](https://via.placeholder.com/15/008000/000000?text=+) Verified
  - which means there are dependencies with no violation.
- ![#FFFF00](https://via.placeholder.com/15/FFFF00/000000?text=+) Ignored
  - which means the path is ignored (not checked).
- ![#800080](https://via.placeholder.com/15/800080/000000?text=+) Warning
  - which means something happened (and gochk didn't check it).
- ![#FF0000](https://via.placeholder.com/15/FF0000/000000?text=+) Violated
  - which means there are dependencies which violates dependency rule.

For `None`, `Verified`, and `Ignored`, only the file path will be displayed.

```
[None]     ../goilerplate/internal/app/adapter/postgresql/conn.go
```

```
[Verified] ../goilerplate/cmd/app/main.go
```

```
[Ignored]  ../goilerplate/.git
```

For `Warning`, it displays what happened to the file.

```
[Warning]  open /Users/resotto/go/src/github.com/resotto/goilerplate/internal/app/application/usecase/lock.go: permission denied
```

For `Violated`, it displays the file path, its dependency, and how it violates dependency rule.

```
[Violated] ../goilerplate/internal/app/domain/temp.go imports "github.com/resotto/goilerplate/internal/app/adapter/postgresql/model"
 => domain depends on adapter
```

## Configuration

`gochk/configs/config.json` has configuration values.

```json
{
  "dependencyOrders": ["external", "adapter", "application", "domain"],
  "ignore": ["test", ".git"],
  "printViolationsAtTheBottom": false
}
```

- `dependencyOrders` are the paths of each circles in Clean Architecture.

  - For example, if you have following four circles, you should specify them from the outer to the inner like: `["external", "adapter", "application", "domain"]`.

    - "External" (most outer)
    - "Adapter"
    - "Application"
    - "Domain" (most inner, the core)

  - If you have other layered architecture, you could specify its layers to this parameter as well.

- `ignore` has the paths ignored by gochk, which can be file path or dir path.

  - If you have the directory you want to ignore, **specifing them might improve the performance of gochk since it returns `filepath.SkipDir`**.

```go
// read.go
func matchIgnore(ignorePaths []string, path string, info os.FileInfo) (bool, error) {
	if included, _ := include(ignorePaths, path); included {
		if info.IsDir() {
			return true, filepath.SkipDir
		}
		return true, nil
	}
	return false, nil
}
```

- `printViolationsAtTheBottom` is the flag whether gochk prints violations of the dependency rule at the bottom or not.

  - If `true`, you can see violations at the bottom like:
  <p align="center">
    <img src="https://user-images.githubusercontent.com/19743841/97001729-d42be180-1573-11eb-90d4-8c68c37f1e04.gif">
  </p>

  - If `false`, you see them disorderly (by goroutine):
  <p align="center">
    <img src="https://user-images.githubusercontent.com/19743841/97001521-74353b00-1573-11eb-8437-fe980c3b34ab.gif">
  </p>

## Unit Testing

Unit test files are located in `gochk/internal/gochk`.

```zsh
gochk
├── internal
│   └── gochk
│       ├── calc_internal_test.go # Unit test (internal)
│       └── read_internal_test.go # Unit test (internal)
└── test
    └── testdata                  # Test data
```

So you can do unit test like:

```zsh
cd internal/gochk
```

```zsh
~/go/src/github.com/resotto/gochk/internal/gochk (master) > go test ./... # Please specify -v if you need detailed outputs
ok      github.com/resotto/gochk/internal/gochk (cached)
```

You can also clean test cache with `go clean -testcache`.

```zsh
~/go/src/github.com/resotto/gochk/internal/gochk (master) > go clean -testcache
~/go/src/github.com/resotto/gochk/internal/gochk (master) > go test ./...
ok      github.com/resotto/gochk/internal/gochk 0.065s # Not cache
```

## Performance Test

Performance test file is located in `gochk/test/performance`.

```zsh
gochk
└── test
    └── performance
        └── check_test.go # Performance test
```

Thus, you can do performance test as follows. It will take few minutes.

```zsh
cd test/performance
```

```zsh
~/go/src/github.com/resotto/gochk/test/performance (master) > go test ./...
ok      github.com/resotto/gochk/test/performance       64.661s
```

### Test Contents

Performance test checks 40,000 test files in `gochk/test/performance` and measures only how long it takes to do it.

#### Note

- Test files will be created before the test and be deleted after the test.
- For each test directory, there will be 10,000 .go test files.

```zsh
gochk
└── test
    ├── performance
    │   ├── adapter         # Test directory
    │   │   ├── postgresql
    │   │   │   └── model
    │   │   ├── repository
    │   │   ├── service
    │   │   ├── view
    │   │   ...             # Test files (g0.go ~ g9999.go)
    │   ├── application     # Test directory
    │   │   ├── service
    │   │   ├── usecase
    │   │   ...             # Test files (g0.go ~ g9999.go)
    │   ├── domain          # Test directory
    │   │   ├── factory
    │   │   ├── repository
    │   │   ├── valueobject
    │   │   ...             # Test files (g0.go ~ g9999.go)
    │   └── external        # Test directory
    │       ...             # Test files (g0.go ~ g9999.go)
    └── testdata
        ├── adapter.txt     # original file of performance/adapter/gX.go
        ├── application.txt # original file of performance/application/gX.go
        ├── domain.txt      # original file of performance/domain/gX.go
        └── external.txt    # original file of performance/external/gX.go
```

For each file, it imports standard libraries and dependencies like:

```go
package xxx

import (
    // standard library imports omitted here

    "github.com/resotto/gochk/test/performance/adapter"                  // import this up to adapter
    "github.com/resotto/gochk/test/performance/adapter/postgresql"       // import this up to adapter
    "github.com/resotto/gochk/test/performance/adapter/postgresql/model" // import this up to adapter
    "github.com/resotto/gochk/test/performance/adapter/repository"       // import this up to adapter
    "github.com/resotto/gochk/test/performance/adapter/service"          // import this up to adapter
    "github.com/resotto/gochk/test/performance/adapter/view"             // import this up to adapter
    "github.com/resotto/gochk/test/performance/application/service"      // import this up to application
    "github.com/resotto/gochk/test/performance/application/usecase"      // import this up to application
    "github.com/resotto/gochk/test/performance/domain/factory"           // import this in only domain
    "github.com/resotto/gochk/test/performance/domain/repository"        // import this in only domain
    "github.com/resotto/gochk/test/performance/domain/valueobject"       // import this in only domain
    "github.com/resotto/gochk/test/performance/external"                 // import this up to adapter
)
```

In performance test, `dependencyOrders` are:

```go
var dependencyOrders = []string{"external", "adapter", "application", "domain"}
```

So, the number of violations equals to:

- domain
  - there are 9 violations x 10,000 files = 90,000
    - domain depends on application (x2)
      ```go
      "github.com/resotto/gochk/test/performance/application/service"
      "github.com/resotto/gochk/test/performance/application/usecase"
      ```
    - domain depends on adapter (x6)
      ```go
      "github.com/resotto/gochk/test/performance/adapter"
      "github.com/resotto/gochk/test/performance/adapter/postgresql"
      "github.com/resotto/gochk/test/performance/adapter/postgresql/model"
      "github.com/resotto/gochk/test/performance/adapter/repository"
      "github.com/resotto/gochk/test/performance/adapter/service"
      "github.com/resotto/gochk/test/performance/adapter/view"
      ```
    - domain depends on external (x1)
      ```go
      "github.com/resotto/gochk/test/performance/external"
      ```
- application
  - there are 7 violations x 10,000 files = 70,000
    - application depends on adapter (x6)
      ```go
      "github.com/resotto/gochk/test/performance/adapter"
      "github.com/resotto/gochk/test/performance/adapter/postgresql"
      "github.com/resotto/gochk/test/performance/adapter/postgresql/model"
      "github.com/resotto/gochk/test/performance/adapter/repository"
      "github.com/resotto/gochk/test/performance/adapter/service"
      "github.com/resotto/gochk/test/performance/adapter/view"
      ```
    - application depends on external (x1)
      ```go
      "github.com/resotto/gochk/test/performance/external"
      ```
- adapter

  - there is 1 violation x 10,000 files = 10,000
    - adapter depends on external (x1)
      ```go
      "github.com/resotto/gochk/test/performance/external"
      ```

- external
  - there are no violations.
- Total
  - 90,000 (domain) + 70,000 (application) + 10,000 (adapter) = 170,000 violations

### Score

Following scores are not cached ones and measured by two Macbook Pros whose spec is different.

| CPU                             | RAM                    | 1st score | 2nd score | 3rd score | Average    |
| :------------------------------ | :--------------------- | :-------- | :-------- | :-------- | :--------- |
| 2.7 GHz Dual-Core Intel Core i5 | 8 GB 1867 MHz DDR3     | 90.43s    | 91.23s    | 88.84s    | **90.17s** |
| 2 GHz Quad-Core Intel Core i5   | 32 GB 3733 MHz LPDDR4X | 54.97s    | 54.03s    | 50.60s    | **53.2s**  |

## Customization

### Changing Result Color

First, please add the ANSI escape code to print.go:

```go
const (
	teal     color = "\033[1;36m"
	green          = "\033[1;32m"
	yellow         = "\033[1;33m"
	purple         = "\033[1;35m"
	red            = "\033[1;31m"
	newColor       = "\033[1;34m" // New color
	reset          = "\033[0m"
)
```

And then, let's change color of result type in read.go:

```go
func newWarning(message string) CheckResult {
	cr := CheckResult{}
	cr.resultType = warning
	cr.message = message
	cr.color = newColor // New color
	return cr
}
```

### Tuning the number of goroutine

If `printViolationsAtTheBottom` is `false`, gochk prints results with goroutine.

You can change the number of goroutine in print.go:

```go
func printConcurrently(results []CheckResult) {
	c := make(chan struct{}, 10) // 10 goroutines by default
	var wg sync.WaitGroup
	for _, r := range results {
		r := r
		c <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() { <-c; wg.Done() }()
			printColorMessage(r)
		}()
	}
	wg.Wait()
}
```

## Build

How to build ...

## Feedback

- [Feel free to write your thoughts](https://github.com/resotto/gochk/issues/1)
- Report a bug to [Bug report](https://github.com/resotto/gochk/issues/2).

## License

[GNU General Public License v3.0](https://github.com/resotto/gochk/blob/master/LICENSE).
