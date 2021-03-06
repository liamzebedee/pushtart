// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autocert

import (
	"crypto"
	"time"

	"golang.org/x/net/context"
)

// maxRandRenew is a maximum deviation from Manager.RenewBefore.
const maxRandRenew = time.Hour

// domainRenewal tracks the state used by the periodic timers
// renewing a single domain's cert.
type domainRenewal struct {
	m      *Manager
	domain string
	key    crypto.Signer
}

// renew is called periodically by a timer.
// The first renew call is kicked off by dr.m.renew.
func (dr *domainRenewal) renew() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	// TODO: rotate dr.key at some point?
	next, err := dr.do(ctx)
	if err != nil {
		next = maxRandRenew / 2
		next += time.Duration(pseudoRand.int63n(int64(next)))
	}
	testDidRenewLoop(next, err)
	time.AfterFunc(next, dr.renew)
}

// do is similar to Manager.createCert but it doesn't lock a Manager.state item.
// Instead, it requests a new certificate independently and, upon success,
// replaces dr.m.state item with a new one and updates cache for the given domain.
//
// It may return immediately if the expiration date of the currently cached cert
// is far enough in the future.
//
// The returned value is a time interval after which the renewal should occur again.
func (dr *domainRenewal) do(ctx context.Context) (time.Duration, error) {
	// a race is likely unavoidable in a distributed environment
	// but we try nonetheless
	if tlscert, err := dr.m.cacheGet(dr.domain); err == nil {
		next := dr.next(tlscert.Leaf.NotAfter)
		if next > dr.m.renewBefore()+maxRandRenew {
			return next, nil
		}
	}

	der, leaf, err := dr.m.authorizedCert(ctx, dr.key, dr.domain)
	if err != nil {
		return 0, err
	}
	state := &certState{
		key:  dr.key,
		cert: der,
		leaf: leaf,
	}
	tlscert, err := state.tlscert()
	if err != nil {
		return 0, err
	}
	dr.m.cachePut(dr.domain, tlscert)
	dr.m.stateMu.Lock()
	defer dr.m.stateMu.Unlock()
	// m.state is guaranteed to be non-nil at this point
	dr.m.state[dr.domain] = state
	return dr.next(leaf.NotAfter), nil
}

func (dr *domainRenewal) next(expiry time.Time) time.Duration {
	d := expiry.Sub(timeNow()) - dr.m.renewBefore()
	// add a bit of randomness to renew deadline
	n := pseudoRand.int63n(int64(maxRandRenew))
	d -= time.Duration(n)
	if d < 0 {
		return 0
	}
	return d
}
