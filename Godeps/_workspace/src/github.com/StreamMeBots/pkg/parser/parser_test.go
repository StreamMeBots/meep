package parser

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/StreamMeBots/pkg/commands"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name    string
		message []byte
		err     bool
		command *commands.Command
	}{
		{
			name:    "pass: SAY command",
			message: []byte(`SAY username="james" message="foo bar baz, message=\"foo bar\" bust"`),
			command: &commands.Command{
				Name: "SAY",
				Args: map[string]string{
					"username": "james",
					"message":  `foo bar baz, message="foo bar" bust`,
				},
			},
		},
		{
			name:    "pass: JOIN command",
			message: []byte(`JOIN username="james" message="foo bar baz, message=\"foo bar\" bust"`),
			command: &commands.Command{
				Name: "JOIN",
				Args: map[string]string{
					"username": "james",
					"message":  `foo bar baz, message="foo bar" bust`,
				},
			},
		},
		{
			name:    "pass: ERROR command",
			message: []byte(`ERROR username="james" message="foo bar baz, message=\"foo bar\" bust"`),
			command: &commands.Command{
				Name: "ERROR",
				Args: map[string]string{
					"username": "james",
					"message":  `foo bar baz, message="foo bar" bust`,
				},
			},
		},
		{
			name:    "pass: SAY command with spaces",
			message: []byte(`SAY username = "james"   message ="foo bar baz\" key=\"foo bar\" bust" publicId= "asdfasdfasdf"`),
			command: &commands.Command{
				Name: "SAY",
				Args: map[string]string{
					"username": "james",
					"message":  `foo bar baz" key="foo bar" bust`,
					"publicId": "asdfasdfasdf",
				},
			},
		},
		{
			name:    "pass: key/value dedupe",
			message: []byte(`JOIN username="james" username="bob" message="foo bar baz\" message=\"foo bar\" bust"`),
			command: &commands.Command{
				Name: "JOIN",
				Args: map[string]string{
					"username": "james",
					"message":  `foo bar baz" message="foo bar" bust`,
				},
			},
		},
		{
			name:    "pass: message with out leading command - test default",
			message: []byte(`username="james" message="foo bar baz, message=\"foo bar\" bust"`),
			command: &commands.Command{
				Name: "SAY",
				Args: map[string]string{
					"username": "james",
					"message":  `foo bar baz, message="foo bar" bust`,
				},
			},
		},
		{
			name:    "fail: invalid command that starts correctly",
			message: []byte(`LEAVEING username="james" message="foo bar baz, message=\"foo bar\" bust"`),
			command: nil,
			err:     true,
		},
		{
			name:    "fail: invalid command with correct starting letter",
			message: []byte(`SAD username="james" message="foo bar baz, message=\"foo bar\" bust"`),
			command: nil,
			err:     true,
		},
		{
			name:    "fail: invalid command with incorrect starting letter",
			message: []byte(`BAD username="james" message="foo bar baz, message=\"foo bar\" bust"`),
			command: nil,
			err:     true,
		},
		{
			name:    "fail: empty message with spaces",
			message: []byte(`       `),
			command: nil,
			err:     true,
		},
		{
			name:    "fail: incomplete command",
			message: []byte(`SAY username="james" message="foo bar baz, message=\"foo bar\" bust`),
			command: nil,
			err:     true,
		},
		{
			name:    "fail: invalid key delimiter",
			message: []byte(`SAY username="james", message="foo bar baz, message=\"foo bar\" bust`),
			command: nil,
			err:     true,
		},
		{
			name:    "fail: invalid key ending equal sign",
			message: []byte(`SAY username+"james", message="foo bar baz, message=\"foo bar\" bust`),
			command: nil,
			err:     true,
		},
		{
			name:    "fail: invalid value start",
			message: []byte(`SAY username=james message="foo bar baz, message=\"foo bar\" bust`),
			command: nil,
			err:     true,
		},
	}

	for _, c := range cases {
		t.Log(c.name)
		cmd, err := ParseBytes(c.message)
		t.Log(err)
		if c.err && err == nil {
			t.Error("	Error should not have been nil but was")
			continue
		}
		if !c.err && err != nil {
			t.Errorf("	Error should of been nil but was not: %v", err.Error())
			continue
		}
		if !reflect.DeepEqual(cmd, c.command) {
			t.Errorf("	commands do not match\n		got:  %+v \n	want: %+v", cmd, c.command)
		}
	}

}

func TestParseReader(t *testing.T) {
	// eof
	eof := struct {
		name    string
		message []byte
		err     bool
		command *commands.Command
	}{
		name:    "fail: EOF",
		message: []byte(`SAY username="james" message="foo bar baz, message=\"foo bar\" bust"`),
		err:     true,
	}
	t.Log(eof.name)
	cmd, err := Parse(bytes.NewReader(eof.message))
	t.Log(err)
	if eof.err && err == nil {
		t.Error("	Error should not have been nil but was")
		return
	}
	if !eof.err && err != nil {
		t.Errorf("	Error should of been nil but was not: %v", err.Error())
		return
	}
	if !reflect.DeepEqual(cmd, eof.command) {
		t.Errorf("	commands do not match\n		got:  %+v \n	want: %+v", cmd, eof.command)
	}
}
