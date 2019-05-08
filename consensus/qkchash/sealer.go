package qkchash

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/QuarkChain/goquarkchain/consensus"
	"github.com/QuarkChain/goquarkchain/core/types"
)

var (
	// two256 is a big integer representing 2^256
	two256 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))

	errNoMiningWork = errors.New("no mining work available yet")
)

func (q *QKCHash) verifySeal(chain consensus.ChainReader, header types.IHeader, adjustedDiff *big.Int) error {
	if header.GetDifficulty().Sign() <= 0 {
		return consensus.ErrInvalidDifficulty
	}
	diff := adjustedDiff
	if diff == nil || diff.Cmp(big.NewInt(0)) == 0 {
		diff = header.GetDifficulty()
	}
	miningRes, err := q.hashAlgo(0 /* TODO: not used */, header.SealHash().Bytes(), header.GetNonce())
	if err != nil {
		return err
	}
	if !bytes.Equal(header.GetMixDigest().Bytes(), miningRes.Digest.Bytes()) {
		return consensus.ErrInvalidMixDigest
	}
	target := new(big.Int).Div(two256, diff)
	if new(big.Int).SetBytes(miningRes.Result).Cmp(target) > 0 {
		return consensus.ErrInvalidPoW
	}
	return nil
}
