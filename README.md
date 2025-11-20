# go-farsi-fake

Generate random Farsi (Persian) fake strings for Go projects, based on farsi hunspell dictionary.

Don't use this package to generate real data for your content because the generated string is absolutely random and may contain undesirable, unpleasant or disagreeable data.

## Installation

To install the package, use the following command:

```sh
go get github.com/saleh-rahimzadeh/go-farsi-fake
```

## Usage

Here's a simple example of how to use the package:

```go
import "github.com/saleh-rahimzadeh/go-farsi-fake"

func main() {
	ff, err := farsifake.New()
	if err != nil {
		panic(err)
	}
	defer ff.Close()

	// Generate a random farsi word.
	str, err := ff.Generate()
	if err != nil {
		panic(err)
	}
	println(str)

	// Generate a sentence with the specified number of farsi words.
	strSnt, err := ff.Sentence(5)
	if err != nil {
		panic(err)
	}
	println(strSnt)

	// Generate a paragraph with a random number of sentences between min and max.
	strPrg, err := ff.Paragraph(3, 7)
	if err != nil {
		panic(err)
	}
	println(strPrg)

	// Generate a range of farsi words as slice.
	slc, err := ff.Range(9)
	if err != nil {
		panic(err)
	}
	println(slc)
}
```

Options:

```go
ff.JumpFromStart = false
ff.BypassError = false
```
