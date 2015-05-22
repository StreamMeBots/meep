package commands

import "testing"

func TestCommand_Get(t *testing.T) {
	var c *Command
	if v := c.Get("foo"); v != "" {
		t.Errorf("Get should of returned: an empty string got: %s", v)
	}

	c = &Command{
		Args: map[string]string{
			"bar": "baz",
		},
	}
	if v := c.Get("bar"); v != "baz" {
		t.Errorf("key 'bar' should of returned: %s", v)
	}
	if v := c.Get("foo"); v != "" {
		t.Errorf("Get should of returned: an empty string got: %s", v)
	}
}

func TestPass(t *testing.T) {
	k := "foo"
	s := "bar"
	r := "PASS " + k + " " + s
	if m := Pass(k, s); m != r {
		t.Errorf("expected: %s got: %s", r, m)
	}
}

func TestSay(t *testing.T) {
	msg := "Hello, World!"
	r := "SAY " + msg
	if m := Say(msg); m != r {
		t.Errorf("expected: %s got: %s", r, m)
	}
}

func TestRoom(t *testing.T) {
	u := "james"
	r := "user:" + u + ":web"
	if m := NewRoom(u); string(m) != r {
		t.Errorf("expected: %s got: %s", r, m)
	}
}
