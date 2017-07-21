// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor: Aaron Meihm ameihm@mozilla.com [:alm]
package main

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/afero"
	"io/ioutil"
	"mig.ninja/mig"
	"path"
	"testing"
	"time"
)

var testRunnerResult = `
{
    "action": {
        "counters": {},
        "description": {},
        "expireafter": "2017-07-10T22:36:00.018410385Z",
        "finishtime": "9998-01-11T11:11:11.000000011Z",
        "id": 8652929898678,
        "lastupdatetime": "0009-01-11T11:11:11.000000011Z",
        "name": "mig-runner: test",
        "operations": [
            {
                "module": "file",
                "parameters": {
                    "searches": {
                        "s1": {
                            "contents": [
                                "testing"
                            ],
                            "names": [
                                "."
                            ],
                            "options": {
                                "macroal": false,
                                "matchall": true,
                                "matchlimit": 1000,
                                "maxdepth": 1000,
                                "maxerrors": 30,
                                "mismatch": null
                            },
                            "paths": [
                                "/home"
                            ],
                            "sizes": [
                                "<1m"
                            ]
                        }
                    }
                }
            }
        ],
        "pgpsignatures": [
            "TESTINGSIGNATURE"
        ],
        "starttime": "0009-01-11T11:11:11.000000011Z",
        "status": "pending",
        "syntaxversion": 2,
        "target": "(tags#>>'{class}') is null",
        "threat": {},
        "validfrom": "2017-07-10T22:34:00.018410385Z"
    },
    "commands": null,
    "name": "test",
    "plugin": "testplugin"
}
`

func TestIncomingResult(t *testing.T) {
	var (
		err error
		res mig.RunnerResult
	)
	ctx, err = initContext("./testdata/runner.cfg")
	if err != nil {
		t.Fatalf("initContext: %v", err)
	}
	ctx.fs = afero.NewMemMapFs()
	go processResults()
	err = json.Unmarshal([]byte(testRunnerResult), &res)
	if err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	ctx.Channels.Results <- res
	extime := time.Now().Add(time.Second * 5)
	// Verify we have the new result queued up in the inflight directory for the entity
	for {
		fd, err := ctx.fs.Open("./testdata/runner/runners/test/inflight/8652929898678.json")
		if err == nil {
			defer fd.Close()
			buf, err := ioutil.ReadAll(fd)
			if err != nil {
				t.Fatalf("ioutil.ReadAll: %v", err)
			}
			comparebuf, err := json.Marshal(res)
			if err != nil {
				t.Fatalf("json.Marshal: %v", err)
			}
			if bytes.Compare(buf, comparebuf) != 0 {
				t.Fatalf("stored inflight action does not equal original")
			}
			break
		}
		if time.Now().After(extime) {
			t.Fatalf("inflight action was not queued in time")
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func TestResultsFetch(t *testing.T) {
	var (
		err error
	)
	ctx, err = initContext("./testdata/runnerfetchresults.cfg")
	if err != nil {
		t.Fatalf("initContext: %v", err)
	}
	go func() {
		for event := range ctx.Channels.Log {
			_, err = mig.ProcessLog(ctx.Logging, event)
			if err != nil {
				t.Fatalf("mig.ProcessLog: %v", err)
			}
		}
	}()
	basefs := afero.NewOsFs()
	roBase := afero.NewReadOnlyFs(basefs)
	ctx.fs = afero.NewCopyOnWriteFs(roBase, afero.NewMemMapFs())
	// Reduce the default timeout delay since we want to try the results fetch
	// immediately
	processResultsTimeout = 0
	go processResults()
	time.Sleep(time.Second * 5)
}

func TestResultsStoragePath(t *testing.T) {
	var err error
	ctx, err = initContext("./testdata/runner.cfg")
	if err != nil {
		t.Fatalf("initContext: %v", err)
	}
	ctx.fs = afero.NewMemMapFs()
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
