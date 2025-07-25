// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eip712

import (
	"fmt"

	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// type aliases to avoid importing "core" everywhere
type TypedData = apitypes.TypedData
type TypedDataDomain = apitypes.TypedDataDomain
type Types = apitypes.Types
type Type = apitypes.Type
type TypedDataMessage = apitypes.TypedDataMessage

// EncodeForSigning encodes the hash that will be signed for the given EIP712 data
func EncodeForSigning(typedData *TypedData) ([]byte, error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, err
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}

	rawData := fmt.Appendf(nil, "\x19\x01%s%s", string(domainSeparator), string(typedDataHash))
	return rawData, nil
}

// EIP712DomainType is the type description for the EIP712 Domain
var EIP712DomainType = []Type{
	{
		Name: "name",
		Type: "string",
	},
	{
		Name: "version",
		Type: "string",
	},
	{
		Name: "chainId",
		Type: "uint256",
	},
}
