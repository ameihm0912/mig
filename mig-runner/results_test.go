// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor: Aaron Meihm ameihm@mozilla.com [:alm]
package main

import (
	"path"
	"testing"
	"time"
)

func TestResultsStoragePath(t *testing.T) {
	var err error
	ctx, err = initContext("./testdata/runner.cfg")
	if err != nil {
		t.Fatalf("initContext: %v", err)
	}
	rdir, err := getResultsStoragePath("test")
	if err != nil {
		t.Fatalf("getResultsStoragePath: %v", err)
	}
	tstamp := time.Now().UTC().Format("20060102")
	comparedir := path.Join("./testdata", "runner", "runners", "test", "results", tstamp)
	if rdir != comparedir {
		t.Fatalf("results storage path %v, expected %v", rdir, comparedir)
	}
}
