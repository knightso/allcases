// want package:`c.TestKind:\[c.TestKindHoge c.TestKindFuga c.TestKindPiyo\] c.TestKind2:\[c.TestKind2Hoge c.TestKind2Fuga c.TestKind2Piyo\]`

package c

type TestKind int

const (
	TestKindHoge TestKind = iota
	TestKindFuga
	TestKindPiyo
)

type TestKind2 int

const (
	TestKind2Hoge TestKind2 = iota
	TestKind2Fuga
	TestKind2Piyo
)

func c() {
	{
		var v TestKind

		// no annotation
		switch v {
		case TestKindHoge:
			// do something
		}

		// allcases
		switch v {
		case TestKindHoge:
		case TestKindFuga:
		case TestKindPiyo:
		default:
			// do something
		}

		// allcases
		switch v { // want "no case of c.TestKindPiyo"
		case TestKindHoge:
		case TestKindFuga:
		default:
			// do something
		}
	}
}
