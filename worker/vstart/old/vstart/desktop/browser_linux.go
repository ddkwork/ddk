// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package desktop provides desktop integration utilities.
package desktop

import (
	"os/exec"

	"github.com/ddkwork/golibrary/std/mylog"
)

// Open asks the system to open the provided path or URL.
func Open(pathOrURL string) {
	mylog.Check(exec.Command("xdg-open", pathOrURL).Start())
	//if err := exec.Command("xdg-open", pathOrURL).Start(); err != nil {
	//	return errs.NewWithCause("Unable to open "+pathOrURL, err)
	//}
}
