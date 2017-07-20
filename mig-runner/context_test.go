// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor: Aaron Meihm ameihm@mozilla.com [:alm]
package main

import (
	"testing"
)

func TestConfigRead(t *testing.T) {
	c, err := initContext("./testdata/runner.cfg")
	if err != nil {
		t.Fatalf("initContext: %v", err)
	}
	if c.Client.Passphrase != "passphrase" {
		t.Fatal("resulting context had invalid passphrase")
	}
}
