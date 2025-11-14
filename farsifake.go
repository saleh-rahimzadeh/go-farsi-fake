package farsifake

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
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
	_JUMP_MIN = 1
	_JUMP_MAX = 356853
)

//┌ Errors
//└─────────────────────────────────────────────────────────────────────────────────────────────────

var (
	ErrFileOpen  = fmt.Errorf("error in opening farsi dictionary file")
	ErrFileClose = fmt.Errorf("error in closing farsi dictionary file")
	ErrGenerate  = fmt.Errorf("error in generating fake farsi word")
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

	file    fs.File
	scanner *bufio.Scanner
	random  *rand.Rand
}

//┌ Internal
//└─────────────────────────────────────────────────────────────────────────────────────────────────

func (o FarsiFake) jump() int {
	return o.random.Intn(_JUMP_MAX-_JUMP_MIN+1) + _JUMP_MIN
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

func (o FarsiFake) Generate() (string, error) {
	if o.JumpFromStart {
		if seeker, ok := o.file.(io.Seeker); ok {
			_, _ = seeker.Seek(0, io.SeekStart)
		}
	}

	counter := o.jump()

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

	if err := o.scanner.Err(); err != nil {
		return "", ErrGenerate
	}

	return "", nil
}
