// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com
package main

import (
	"fmt"
	"mig"
)

func injectMessage(msg string, rkey string, ctx Context) {
	ctx.Channels.Log <- mig.Log{Desc: "entering injectMessage()"}.Debug()

	ctx.Channels.Log <- mig.Log{Desc: fmt.Sprintf("injecting to %v '%v'", rkey, msg)}.Debug()

	publication.Unlock()
	defer publication.Lock()

	err := publish(ctx, "mig", rkey, []byte(msg))
	if err != nil {
		panic(err)
	}
}
