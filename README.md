# allcases

This Repository is forked from https://github.com/knightso/allcases

`allcalses` is Go code analyzer which checks for switch statements to have all cases.

# Install
```
$ go get github.com/knightso/allcases/cmd/allcases
```

# How to use

## Usage

```
$ allcases [-flag] [package]
```

- `flags` are only ones took over from https://godoc.org/golang.org/x/tools/go/analysis tool. check `allcases -help` for them.
- `package` format is based on go tools such as `./...` .

## Annotation

This checker checks for switch statements which have an annotation comment `// allcases`
to have all cases for consts of the expression type.

For example if there is a type and consts like:

```go
	type TestKind int

	const (
		TestKindHoge TestKind = iota
		TestKindFuga
		TestKindPiyo
	)
```

and switch statements like:

```go
	// allcases
	switch v {
	case TestKindHoge:
		// do something
	case TestKindFuga:
		// do something
	}
```

then the checker reports that the switch statement doesn't have the TestKindPiyo case.
