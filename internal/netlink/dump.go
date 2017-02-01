// Copyright © 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package netlink

import (
	"fmt"
	"io"

	"github.com/platinasystems/go/internal/indent"
)

var usage = `
usage: nldump [-all-nsid] [ noop|error|done|link|addr|route|neighor ]...
`[1:]

// Dump all or the selected netlink messages
func Dump(w io.Writer, args ...string) error {
	var allnsid,
		mayDumpNoop,
		mayDumpError,
		mayDumpDone,
		mayDumpLink,
		mayDumpAddr,
		mayDumpRoute,
		mayDumpNeighbor bool
	for _, arg := range args {
		switch arg {
		case "-h", "-help", "--help":
			fmt.Fprint(w, usage)
			return nil
		case "-all-nsid", "--all-nsid", "all-nsid":
			allnsid = true
		case "noop":
			mayDumpNoop = true
		case "error":
			mayDumpError = true
		case "done":
			mayDumpDone = true
		case "link":
			mayDumpLink = true
		case "addr":
			mayDumpAddr = true
		case "route":
			mayDumpRoute = true
		case "neighbor":
			mayDumpNeighbor = true
		default:
			return fmt.Errorf("%s: unknown", arg)
		}
	}
	if !mayDumpNoop && !mayDumpError && !mayDumpDone && !mayDumpLink &&
		!mayDumpAddr && !mayDumpRoute && !mayDumpNeighbor {

		mayDumpNoop = true
		mayDumpError = true
		mayDumpDone = true
		mayDumpLink = true
		mayDumpAddr = true
		mayDumpRoute = true
		mayDumpNeighbor = true
	}
	rx := make(chan Message, 64)
	nl, err := New(rx)
	if err != nil {
		return err
	}
	if allnsid {
		err = nl.ListenAllNsid()
		if err != nil {
			return err
		}
	}
	go nl.Listen()
	for msg := range rx {
		_, isNoop := msg.(*NoopMessage)
		_, isError := msg.(*ErrorMessage)
		_, isDone := msg.(*DoneMessage)
		_, isLink := msg.(*IfInfoMessage)
		_, isAddr := msg.(*IfAddrMessage)
		_, isRoute := msg.(*RouteMessage)
		_, isNeighbor := msg.(*NeighborMessage)
		if (mayDumpNoop && isNoop) ||
			(mayDumpError && isError) ||
			(mayDumpDone && isDone) ||
			(mayDumpLink && isLink) ||
			(mayDumpAddr && isAddr) ||
			(mayDumpRoute && isRoute) ||
			(mayDumpNeighbor && isNeighbor) {
			_, err = msg.WriteTo(indent.New(w, "    "))
			if err != nil {
				return err
			}
		}
		msg.Close()
	}
	return nil
}
