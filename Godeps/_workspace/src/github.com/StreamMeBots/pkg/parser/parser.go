/*
* Package parser is used to parse content from stream.me chat servers
*
* Examples Formats:
*
*	JOIN username="joe"
*	SAY username="joe" message="what's up"
 */
package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode"

	"github.com/StreamMeBots/pkg/commands"
)

// Errors
var (
	ErrIncompleteMessage = errors.New("parser: incomplete message")
	ErrEmptyMessage      = errors.New("parser: empty message")
)

// states
const (
	_ = iota
	commandState
	keyState
	valueState
	skipState
)

// ParseBytes is a thin wrapper around Parse
func ParseBytes(msg []byte) (*commands.Command, error) {
	return Parse(bytes.NewReader(append(msg, '\n')))
}

// Parse parses the input read from the Reader to create a Command
//
// NOTE: delim = \n
func Parse(r io.Reader) (*commands.Command, error) {
	rd := bufio.NewReader(r)
	s := &scanner{
		next:       commandStart,
		currentMsg: []byte{},
		command:    &commands.Command{Args: map[string]string{}},
		r:          rd,
	}

	i := 0
	for {
		b, err := rd.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == '\n' {
			break
		}
		r := rune(b)
		s.currentMsg = append(s.currentMsg, b)
		if err := s.next(s, r); err != nil {
			return nil, err
		}
		s.last = r

		i++
	}
	if len(s.state) == 0 && len(s.command.Name) != 0 {
		return s.command, nil
	}

	return nil, fmt.Errorf("parser: incomplete command")
}

// scanner is used to parse the message
type scanner struct {
	command      *commands.Command
	state        []int
	next         func(*scanner, rune) error
	last         rune
	commandState []rune
	currentMsg   []byte
	currentArg   arg
	r            *bufio.Reader
}

// arg represents a key value argument
type arg struct {
	key   string
	value string
}

// pushState appends to the scanners state
func (s *scanner) pushState(state int) {
	s.state = append(s.state, state)
}

// popState pops state from the scanners state
func (s *scanner) popState() int {
	var state int
	s.state, state = s.state[:len(s.state)-1], s.state[len(s.state)-1]
	return state
}

// currentMessage reads the currently parsed message
func (s *scanner) currentMessage() string {
	return string(s.currentMsg)
}

// addArg adds the scanners currentArg to the scanners command Args map
func (s *scanner) addArg() {
	// dropped duped key/value
	if _, ok := s.command.Args[s.currentArg.key]; !ok {
		s.command.Args[s.currentArg.key] = s.currentArg.value
	}
	s.currentArg.key = ""
	s.currentArg.value = ""
}

func commandStart(s *scanner, r rune) error {
	if unicode.IsSpace(r) {
		return nil
	}
	if unicode.IsLower(r) {
		if err := s.r.UnreadByte(); err != nil {
			return err
		}
		s.next = keyStart
		s.command.Name = "SAY"
		return nil
	}

	switch r {
	case 'J': // JOIN
		s.commandState = []rune(`OIN`)
	case 'L': // LEAVE
		s.commandState = []rune(`EAVE`)
	case 'S': // SAY
		s.commandState = []rune(`AY`)
	case 'E': // ERROR
		s.commandState = []rune(`RROR`)
	case 'P': // PASS
		s.commandState = []rune(`ASS`)
	default:
		return fmt.Errorf("parser: invalid command")
	}

	s.next = commandEnd
	s.pushState(commandState)

	return nil
}

func commandEnd(s *scanner, r rune) error {
	if len(s.commandState) == 0 {
		s.command.Name = string(s.currentMsg[:len(s.currentMsg)-1])
		s.popState()
		s.next = keyStart
		return nil
	}
	if s.commandState[0] == r {
		s.commandState = s.commandState[1:]
		return nil
	}

	return fmt.Errorf("parser: invalid command: %v", s.currentMessage())
}

func keyStart(s *scanner, r rune) error {
	if unicode.IsSpace(r) {
		return nil
	}

	if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
		return fmt.Errorf("parser: invalid key delimiter: '%s' - message: %s", string(r), s.currentMessage())
	}

	s.pushState(keyState)
	s.currentArg.key = string(r)
	s.next = keyEnd

	return nil
}

func keyEnd(s *scanner, r rune) error {
	if unicode.IsLetter(r) || unicode.IsNumber(r) {
		if unicode.IsSpace(s.last) {
			return fmt.Errorf("parser: invalid key ending: '%s' - message: %s", string(r), s.currentMessage())
		}
		s.currentArg.key += string(r)
		return nil
	}

	if unicode.IsSpace(r) {
		return nil
	}

	if r == '=' {
		s.next = valueStart
		return nil
	}

	return fmt.Errorf("parser: invalid key ending, looking for '=' got: '%s' - message: %s", string(r), s.currentMessage())
}

func valueStart(s *scanner, r rune) error {
	if r == '"' {
		s.popState()
		s.pushState(valueState)
		s.next = valueEnd
		return nil
	}
	if unicode.IsSpace(r) {
		return nil
	}

	return fmt.Errorf("parser: invalid start value character: %s - message: %s", string(r), s.currentMessage())
}

func valueEnd(s *scanner, r rune) error {
	if s.state[len(s.state)-1] == skipState {
		s.currentArg.value += string(r)
		s.popState()
		return nil
	}
	if r == '\\' {
		s.pushState(skipState)
		return nil
	}
	if r == '"' {
		s.popState()
		s.next = keyStart
		s.addArg()
		return nil
	}

	s.currentArg.value += string(r)

	return nil
}
