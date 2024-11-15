<div align="center">

# comver

</div>

<div align="center">

[![Go](https://github.com/typisttech/comver/actions/workflows/go.yml/badge.svg)](https://github.com/typisttech/comver/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/typisttech/comver/graph/badge.svg?token=GVO7RV80TJ)](https://codecov.io/gh/typisttech/comver)
[![Go Report Card](https://goreportcard.com/badge/github.com/typisttech/comver)](https://goreportcard.com/report/github.com/typisttech/comver)
[![GitHub Release](https://img.shields.io/github/v/release/typisttech/comver?style=flat-square&)](https://github.com/typisttech/comver/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/typisttech/comver.svg)](https://pkg.go.dev/github.com/typisttech/comver)
[![license](https://img.shields.io/github/license/typisttech/comver.svg?style=flat-square)](https://github.com/typisttech/comver/blob/master/LICENSE)
[![X Follow @TangRufus](https://img.shields.io/badge/Follow-%40TangRufus-black?style=flat-square&logo=x&logoColor=white)](https://x.com/tangrufus)
[![Hire Typist Tech](https://img.shields.io/badge/Hire-Typist%20Tech-ff69b4.svg?style=flat-square)](https://typist.tech/contact/)

</div>

<p align="center">
  <strong>Package <code>comver</code> provides the ability to work with <a href="https://github.com/composer/semver/">composer supported versions</a> in Go.</strong>
  <br />
  <br />
  Built with ♥ by <a href="https://typist.tech/">Typist Tech</a>
</p>

---

## Usage

> [!NOTE]
> See full API documentation at [pkg.go.dev](https://pkg.go.dev/github.com/typisttech/comver).

### `Version`

[`NewVersion`](https://pkg.go.dev/github.com/typisttech/comver#NewVersion) parses a given version string, attempts to coerce a version string into a [`Version`](https://pkg.go.dev/github.com/typisttech/comver#Version) object or return an error if unable to parse the version string.

If there is a leading **v** or a version listed without all parts (e.g. **v1.2.p5+foo**) it will attempt to coerce it into a valid composer version (e.g. **1.2.0.0-patch5**). In both cases a [`Version`](https://pkg.go.dev/github.com/typisttech/comver#Version) object is returned that can be sorted, compared, and used in constraints.


> [!WARNING]
> Due to implementation complexity, it only supports a subset of [composer versioning](https://github.com/composer/semver/). 
>
> Refer to the [`version_test.go`](version_test.go) for examples.


```go
ss := []string{
    "1.2.3",
    "v1.2.p5+foo",
    "v1.2.3.4.p5+foo",
    "2010-01-02",
    "2010-01-02.5",
    "not a version",
    "1.0.0-meh",
    "20100102.0.3.4",
    "1.0.0-alpha.beta",
}

for _, s := range ss {
    v, err := comver.NewVersion(s)
    if err != nil {
        fmt.Println(s, " => ", err)
        continue
    }
    fmt.Println(s, " => ", v)
}

// Output:
// 1.2.3  =>  1.2.3.0
// v1.2.p5+foo  =>  1.2.0.0-patch5
// v1.2.3.4.p5+foo  =>  1.2.3.4-patch5
// 2010-01-02  =>  2010.1.2.0
// 2010-01-02.5  =>  2010.1.2.5
// not a version  =>  error parsing version string "not a version"
// 1.0.0-meh  =>  error parsing version string "1.0.0-meh"
// 20100102.0.3.4  =>  error parsing version string "20100102.0.3.4"
// 1.0.0-alpha.beta  =>  error parsing version string "1.0.0-alpha.beta"
```

### `constraint`

```go
v1, _ := comver.NewVersion("1")
v2, _ := comver.NewVersion("2")
v3, _ := comver.NewVersion("3")
v4, _ := comver.NewVersion("4")

cs := []any{
    comver.NewGreaterThanConstraint(v1),
    comver.NewGreaterThanOrEqualToConstraint(v2),
    comver.NewLessThanOrEqualToConstraint(v3),
    comver.NewLessThanConstraint(v4),
}

for _, c := range cs {
    fmt.Println(c)
}

// Output:
// >1
// >=2
// <=3
// <4
```

### `interval`

`interval` represents the intersection (logical AND) of two constraints.

```go
v1, _ := comver.NewVersion("1")
v2, _ := comver.NewVersion("2")
v3, _ := comver.NewVersion("3")

g1l3, _ := comver.NewInterval(
    comver.NewGreaterThanConstraint(v1),
    comver.NewLessThanConstraint(v3),
)

if g1l3.Check(v2) {
    fmt.Println(v2.Short(), "satisfies", g1l3)
}

if !g1l3.Check(v3) {
    fmt.Println(v2.Short(), "doesn't satisfy", g1l3)
}

// Output:
// 2 satisfies >1 <3
// 2 doesn't satisfy >1 <3
```

### `Intervals`

[`Intervals`](https://pkg.go.dev/github.com/typisttech/comver#Intervals) represent the union (logical OR) of multiple intervals.

```go
v1, _ := comver.NewVersion("1")
v2, _ := comver.NewVersion("2")
v3, _ := comver.NewVersion("3")
v4, _ := comver.NewVersion("4")

g1l3, _ := comver.NewInterval(
comver.NewGreaterThanConstraint(v1),
comver.NewLessThanConstraint(v3),
)

ge2le4, _ := comver.NewInterval(
comver.NewGreaterThanOrEqualToConstraint(v2),
comver.NewLessThanOrEqualToConstraint(v4),
)

is := comver.Intervals{g1l3, ge2le4}
fmt.Println(is)

is = comver.Compact(is)
fmt.Println(is)

// Output:
// >1 <3 || >=2 <=4
// >1 <=4
```

## Credits

[`comver`](https://github.com/typisttech/comver) is a [Typist Tech](https://typist.tech) project and maintained by [Tang Rufus](https://x.com/TangRufus), freelance developer for [hire](https://typist.tech/contact/).

Full list of contributors can be found [here](https://github.com/typisttech/comver/graphs/contributors).

## Copyright and License

This project is a [free software](https://www.gnu.org/philosophy/free-sw.en.html) distributed under the terms of the MIT license. For the full license, see [LICENSE](./LICENSE).

## Contribute

Feedbacks / bug reports / pull requests are welcome.
