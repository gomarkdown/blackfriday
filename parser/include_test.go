package parser

import "testing"

func TestUpdateWd(t *testing.T) {
	p := "/tmp"

	if x := updateWd(p, "/new/foo"); x != "/new" {
		t.Errorf("want %s, got %s", "/new", x)
	}

	if x := updateWd(p, "new/new"); x != "/tmp/new" {
		t.Errorf("want %s, got %s", "/tmp/new", x)
	}
}

func TestIsInclude(t *testing.T) {
	p := New()
	if name, _ := p.isInclude([]byte("{{foo}}")); name != "foo" {
		t.Errorf("want %s, got %s", "foo", name)
	}
	if name, _ := p.isInclude([]byte("{{foo}")); name != "" {
		t.Errorf("want %s, got %s", "", name)
	}
	if name, _ := p.isInclude([]byte("{foo}")); name != "" {
		t.Errorf("want %s, got %s", "", name)
	}
}
