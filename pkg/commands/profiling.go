// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/spf13/pflag"
)

var (
	profileName   string
	profileOut string
)

func addProfilingFlags(flags *pflag.FlagSet) {
	flags.StringVar(&profileName, "profile", "none", "Name of profile to capture. One of (none|cpu|heap|goroutine|threadcreate|block|mutex)")
	flags.StringVar(&profileOut, "profile-out", "profile.pprof", "Name of the file to write the profile to")
	flags.MarkHidden("profile")
	flags.MarkHidden("profile-out")
}

func initProfiling() error {
	switch profileName {
	case "none":
		return nil
	case "cpu":
		f, err := os.Create(profileOut)
		if err != nil {
			return err
		}
		return pprof.StartCPUProfile(f)
	// Block and mutex profiles need a call to Set{Block,Mutex}ProfileRate to
	// output anything. We choose to sample all events.
	case "block":
		runtime.SetBlockProfileRate(1)
		return nil
	case "mutex":
		runtime.SetMutexProfileFraction(1)
		return nil
	default:
		// Check the profile name is valid.
		if profile := pprof.Lookup(profileName); profile == nil {
			return fmt.Errorf("unknown profile '%s'", profileName)
		}
	}

	return nil
}

func flushProfiling() error {
	switch profileName {
	case "none":
		return nil
	case "cpu":
		pprof.StopCPUProfile()
	case "heap":
		runtime.GC()
		fallthrough
	default:
		profile := pprof.Lookup(profileName)
		if profile == nil {
			return nil
		}
		f, err := os.Create(profileOut)
		if err != nil {
			return err
		}
		profile.WriteTo(f, 0)
	}

	return nil
}
