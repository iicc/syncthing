// Copyright (C) 2016 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"github.com/syncthing/syncthing/lib/config"
	"math/rand"
	"time"
)

type rescanRequest struct {
	subdirs []string
	err     chan error
}

// bundle all folder scan activity
type folderScanner struct {
	interval time.Duration
	timer    *time.Timer
	now      chan rescanRequest
	delay    chan time.Duration
}

func newFolderScanner(config config.FolderConfiguration) folderScanner {
	return folderScanner{
		interval: time.Duration(config.RescanIntervalS) * time.Second,
		timer:    time.NewTimer(time.Millisecond), // The first scan should be done immediately.
		now:      make(chan rescanRequest),
		delay:    make(chan time.Duration),
	}
}

func (s *folderScanner) reschedule() {
	if s.interval == 0 {
		return
	}
	// Sleep a random time between 3/4 and 5/4 of the configured interval.
	sleepNanos := (s.interval.Nanoseconds()*3 + rand.Int63n(2*s.interval.Nanoseconds())) / 4
	interval := time.Duration(sleepNanos) * time.Nanosecond
	l.Debugln(s, "next rescan in", interval)
	s.timer.Reset(interval)
}

func (s *folderScanner) Scan(subdirs []string) error {
	req := rescanRequest{
		subdirs: subdirs,
		err:     make(chan error),
	}
	s.now <- req
	return <-req.err
}

func (s *folderScanner) Delay(next time.Duration) {
	s.delay <- next
}
