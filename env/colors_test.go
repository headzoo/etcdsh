package env

import (
	"os"
	"testing"
)

func TestGetLS(t *testing.T) {
	os.Setenv("LS_COLORS", "di=01;31:fi=0:ln=31:pi=5:so=5:bd=5:cd=5:or=31:*.deb=90")
	c := NewColors()
	actual, err := c.GetLS("ln")
	if "31" != actual {
		t.Errorf("GetLS('ln') = %v, want '31'.", actual)
	} else if err != nil {
		t.Error(err)
	}

	actual, err = c.GetLS("di")
	if "01;31" != actual {
		t.Errorf("GetLS('di') = %v, want '01;31'.", actual)
	} else if err != nil {
		t.Error(err)
	}

	actual, err = c.GetLS("foo")
	if err == nil {
		t.Error("GetLS('foo') expected error, got nil.")
	}
}

func TestGetLSDefault(t *testing.T) {
	os.Setenv("LS_COLORS", "di=01;31:fi=0:ln=31:pi=5:so=5:bd=5:cd=5:or=31:*.deb=90")
	c := NewColors()
	actual, _ := c.GetLSDefault("ln", "45")
	if "31" != actual {
		t.Errorf("GetLS('ln') = %v, want '31'.", actual)
	}

	actual, _ = c.GetLSDefault("foo", "42")
	if "42" != actual {
		t.Errorf("GetLS('foo') = %v, want '42'.", actual)
	}
}
