package key

import "testing"

func TestIsContext(t *testing.T) {
	if !IsContext(ContextUser) {
		t.Fail()
	}

	if !IsContext(ContextRedirect) {
		t.Fail()
	}

	if IsContext("User") {
		t.Fail()
	}

	if IsContext("RedirectTo") {
		t.Fail()
	}
}
