package farsifake

import (
	"bufio"
	"embed"
	"errors"
	"io"
	"io/fs"
	"math/rand"
	"strings"
	"time"
)

//┌ Embed
//└─────────────────────────────────────────────────────────────────────────────────────────────────

//go:embed fa.dic
var embedFarsiDic embed.FS

//┌ Definitions
//└─────────────────────────────────────────────────────────────────────────────────────────────────

// Farsi dictionary file name.
const _FILE string = "fa.dic"

// Jump range based on number of lines in farsi dictionary file.
const (
	_JUMP_MIN int = 1
	_JUMP_MAX int = 356853
)

//┌ Errors
//└─────────────────────────────────────────────────────────────────────────────────────────────────

var (
	ErrFileOpen      = errors.New("error in opening farsi dictionary file")
	ErrFileClose     = errors.New("error in closing farsi dictionary file")
	ErrGenerate      = errors.New("error in generating fake farsi word")
	ErrInvalidMinMax = errors.New("invalid min and max values")
	ErrInvalidCount  = errors.New("invalid count value")
)

//┌ Instance
//└─────────────────────────────────────────────────────────────────────────────────────────────────

func New() (FarsiFake, error) {
	file, err := embedFarsiDic.Open(_FILE)
	if err != nil {
		return FarsiFake{}, ErrFileOpen
	}

	return FarsiFake{
		JumpFromStart: false,

		file:    file,
		scanner: bufio.NewScanner(file),
		random:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

//┌ Object
//└─────────────────────────────────────────────────────────────────────────────────────────────────

type FarsiFake struct {
	// If true, each call starts from the beginning of the dictionary, else scanning continues from the last position.
	//
	// Default is false to provide better random result.
	JumpFromStart bool

	// If true, errors are bypassed and empty results are returned instead.
	//
	// Default is false to enforce error handling.
	BypassError bool

	file    fs.File
	scanner *bufio.Scanner
	random  *rand.Rand
}

//┌ Internal
//└─────────────────────────────────────────────────────────────────────────────────────────────────

func (o FarsiFake) jump(min int, max int) int {
	return o.random.Intn(max-min+1) + min
}

//┌ Public
//└─────────────────────────────────────────────────────────────────────────────────────────────────

func (o FarsiFake) Close() error {
	err := o.file.Close()
	if err != nil {
		return ErrFileClose
	}
	return nil
}

// Generate a random farsi word.
func (o FarsiFake) Generate() (string, error) {
	if o.JumpFromStart {
		if seeker, ok := o.file.(io.Seeker); ok {
			_, _ = seeker.Seek(0, io.SeekStart)
		}
	}

	counter := o.jump(_JUMP_MIN, _JUMP_MAX)

	var hasNext bool = true
	for hasNext {
		hasNext = o.scanner.Scan()
		if !hasNext {
			if seeker, ok := o.file.(io.Seeker); ok {
				_, _ = seeker.Seek(0, io.SeekStart)
			}
			o.scanner = bufio.NewScanner(o.file)
			hasNext = true
			continue
		}

		counter--
		if counter == 0 {
			return o.scanner.Text(), nil
		}
	}

	if err := o.scanner.Err(); err != nil && !o.BypassError {
		return "", ErrGenerate
	}

	return "", nil
}

// Generate a sentence with the specified number of farsi words.
func (o FarsiFake) Sentence(count int) (string, error) {
	if count < 1 {
		return "", ErrInvalidCount
	}

	sentences := make([]string, 0, count)

	for range count {
		word, err := o.Generate()
		if err != nil {
			return "", ErrGenerate
		}
		sentences = append(sentences, word)
	}

	return strings.Join(sentences, " "), nil
}

// Generate a paragraph with a random number of sentences between min and max.
func (o FarsiFake) Paragraph(min int, max int) (string, error) {
	if min < 1 || max < 1 {
		return "", ErrInvalidCount
	}
	if max < min {
		return "", ErrInvalidMinMax
	}

	count := o.jump(min, max)

	word, err := o.Sentence(count)
	if err != nil {
		return "", ErrGenerate
	}

	return word, nil
}
