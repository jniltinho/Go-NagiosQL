package nagimport

import (
	"os"
	"testing"
)

func writeTmpCfg(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "nagios-*.cfg")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestParseHost(t *testing.T) {
	path := writeTmpCfg(t, `
# test file
define host {
    host_name   web01
    alias       Web Server 01
    address     10.0.0.1
    check_command check-host-alive
    register    1
}
`)
	objs, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(objs) != 1 {
		t.Fatalf("expected 1 object, got %d", len(objs))
	}
	obj := objs[0]
	if obj.Type != "host" {
		t.Errorf("expected type 'host', got %q", obj.Type)
	}
	if obj.Fields["host_name"] != "web01" {
		t.Errorf("expected host_name=web01, got %q", obj.Fields["host_name"])
	}
	if obj.Fields["address"] != "10.0.0.1" {
		t.Errorf("expected address=10.0.0.1, got %q", obj.Fields["address"])
	}
}

func TestParseMultipleObjects(t *testing.T) {
	path := writeTmpCfg(t, `
define host {
    host_name   host1
    address     1.1.1.1
}

define service {
    service_description   PING
    check_command         check_ping!100,20%!500,60%
}

define command {
    command_name   check_ssh
    command_line   $USER1$/check_ssh $ARG1$ $HOSTADDRESS$
}
`)
	objs, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(objs) != 3 {
		t.Fatalf("expected 3 objects, got %d", len(objs))
	}
	types := map[string]bool{}
	for _, o := range objs {
		types[o.Type] = true
	}
	for _, want := range []string{"host", "service", "command"} {
		if !types[want] {
			t.Errorf("missing object type %q", want)
		}
	}
}

func TestParseInlineComment(t *testing.T) {
	path := writeTmpCfg(t, `
define host {
    host_name   myhost ; this is a comment
    address     2.2.2.2
}
`)
	objs, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if objs[0].Fields["host_name"] != "myhost" {
		t.Errorf("inline comment not stripped: %q", objs[0].Fields["host_name"])
	}
}

func TestParseEmptyFile(t *testing.T) {
	path := writeTmpCfg(t, "# only comments\n")
	objs, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(objs) != 0 {
		t.Errorf("expected 0 objects, got %d", len(objs))
	}
}
