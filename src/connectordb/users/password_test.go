package users

/**
This file provides testing for the password management functions

Copyright 2015 - Joseph Lewis <joseph@josephlewis.net>
                 Daniel Kumor <rdkumor@gmail.com>

All Rights Reserved
**/

import (
	"testing"
)


func TestCalcHash(t *testing.T){
	h1 := calcHash("password", "", "")
	h2 := calcHash("password", "", "SHA512")
	h3 := calcHash("password", "a", "SHA512")
	h4 := calcHash("password2", "a", "SHA512")
	h5 := calcHash("password2", "a", "SHA512")

	if h1 != h2 {
		t.Errorf("h1 and h2 should match")
	}

	if h2 == h3 {
		t.Errorf("h2 and h3 should not match")
	}

	if h3 == h4 {
		t.Errorf("h4 and h3 should not match")
	}

	if h5 != h4 {
		t.Errorf("h5 and h4 should match")
	}
}
