// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transaction_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethersphere/bee/v2/pkg/transaction"
	"github.com/ethersphere/bee/v2/pkg/transaction/backendmock"
)

func TestIsSynced(t *testing.T) {
	t.Parallel()

	maxDelay := 10 * time.Second
	now := time.Now().UTC()
	ctx := context.Background()

	t.Run("synced", func(t *testing.T) {
		t.Parallel()

		synced, _, err := transaction.IsSynced(
			ctx,
			backendmock.New(
				backendmock.WithHeaderbyNumberFunc(func(ctx context.Context, number *big.Int) (*types.Header, error) {
					if number != nil {
						return nil, errors.New("latest block should be called with nil")
					}
					return &types.Header{
						Time: uint64(now.Unix()),
					}, nil
				}),
			),
			maxDelay,
		)
		if err != nil {
			t.Fatal(err)
		}
		if !synced {
			t.Fatal("expected synced")
		}
	})

	t.Run("not synced", func(t *testing.T) {
		t.Parallel()

		synced, _, err := transaction.IsSynced(
			ctx,
			backendmock.New(
				backendmock.WithHeaderbyNumberFunc(func(ctx context.Context, number *big.Int) (*types.Header, error) {
					if number != nil {
						return nil, errors.New("latest block should be called with nil")
					}
					return &types.Header{
						Time: uint64(now.Add(-maxDelay).Unix()),
					}, nil
				}),
			),
			maxDelay,
		)
		if err != nil {
			t.Fatal(err)
		}
		if synced {
			t.Fatal("expected not synced")
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("err")
		_, _, err := transaction.IsSynced(
			ctx,
			backendmock.New(
				backendmock.WithHeaderbyNumberFunc(func(ctx context.Context, number *big.Int) (*types.Header, error) {
					if number != nil {
						return nil, errors.New("latest block should be called with nil")
					}
					return nil, expectedErr
				}),
			),
			maxDelay,
		)
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error. wanted %v, got %v", expectedErr, err)
		}
	})
}
