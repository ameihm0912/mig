// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor: Aaron Meihm ameihm@mozilla.com [:alm]
package main

import (
	"testing"
)

func TestLoadPlugins(t *testing.T) {
	var err error
	ctx, err = initContext("./testdata/runner.cfg")
	if err != nil {
		t.Fatalf("initContext: %v", err)
	}
	err = loadPlugins()
	if err != nil {
		t.Fatalf("loadPlugins: %v", err)
	}
	// We should only have the valid executable plugins in the
	// plugin list at this point
	if len(pluginList) != 1 {
		t.Fatalf("plugin count expected %v, have %v", "1", len(pluginList))
	}
}
