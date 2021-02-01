// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/ethersphere/bee/pkg/jsonhttp"
	"github.com/ethersphere/bee/pkg/storage"
	"github.com/ethersphere/bee/pkg/swarm"
	"github.com/ethersphere/bee/pkg/traversal"
	"github.com/gorilla/mux"
)

// pinFile is used to pin an already uploaded content.
func (s *server) pinPin(w http.ResponseWriter, r *http.Request) {
	addr, err := swarm.ParseHexAddress(mux.Vars(r)["address"])
	if err != nil {
		s.Logger.Debugf("pin: parse address: %v", err)
		s.Logger.Error("pin: parse address")
		jsonhttp.BadRequest(w, "bad address")
		return
	}

	has, err := s.Storer.Has(r.Context(), addr)
	if err != nil {
		s.Logger.Debugf("pin: localstore has: %v", err)
		s.Logger.Error("pin: store")
		jsonhttp.InternalServerError(w, err)
		return
	}

	if !has {
		_, err := s.Storer.Get(r.Context(), storage.ModeGetRequest, addr)
		if err != nil {
			s.Logger.Debugf("pin chunk: netstore get: %v", err)
			s.Logger.Error("pin chunk: netstore")
			jsonhttp.NotFound(w, nil)
			return
		}
	}

	ctx := r.Context()

	chunkAddressFn := s.pinChunkAddressFn(ctx, addr)

	err = s.Traversal.TraverseAddresses(ctx, addr, chunkAddressFn)
	if err != nil {
		s.Logger.Debugf("pin: traverse chunks: %v, addr %s", err, addr)

		if errors.Is(err, traversal.ErrInvalidType) {
			s.Logger.Error("pin: invalid type")
			jsonhttp.BadRequest(w, "invalid type")
			return
		}

		s.Logger.Error("pin: cannot pin")
		jsonhttp.InternalServerError(w, "cannot pin")
		return
	}

	jsonhttp.OK(w, nil)
}

// unpinFile removes pinning from content.
func (s *server) unpinPin(w http.ResponseWriter, r *http.Request) {
	addr, err := swarm.ParseHexAddress(mux.Vars(r)["address"])
	if err != nil {
		s.Logger.Debugf("pin: parse address: %v", err)
		s.Logger.Error("pin: parse address")
		jsonhttp.BadRequest(w, "bad address")
		return
	}

	has, err := s.Storer.Has(r.Context(), addr)
	if err != nil {
		s.Logger.Debugf("pin: localstore has: %v", err)
		s.Logger.Error("pin: store")
		jsonhttp.InternalServerError(w, err)
		return
	}

	if !has {
		jsonhttp.NotFound(w, nil)
		return
	}

	ctx := r.Context()

	chunkAddressFn := s.unpinChunkAddressFn(ctx, addr)

	err = s.Traversal.TraverseAddresses(ctx, addr, chunkAddressFn)
	if err != nil {
		s.Logger.Debugf("pin: traverse chunks: %v, addr %s", err, addr)
		s.Logger.Error("pin: cannot unpin")
		jsonhttp.InternalServerError(w, "cannot unpin")
		return
	}

	jsonhttp.OK(w, nil)
}

// getPin shows info on a pinned content root address
func (s *server) getPin(w http.ResponseWriter, r *http.Request) {
	addr, err := swarm.ParseHexAddress(mux.Vars(r)["address"])
	if err != nil {
		s.Logger.Debugf("pin: parse address: %v", err)
		s.Logger.Error("pin: parse address")
		jsonhttp.BadRequest(w, "bad address")
		return
	}
}

// listPins lists pinned content root addresses
func (s *server) listPins(w http.ResponseWriter, r *http.Request) {
}
