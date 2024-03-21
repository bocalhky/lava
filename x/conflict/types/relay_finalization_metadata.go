package types

import (
	"encoding/json"
	fmt "fmt"

	tendermintcrypto "github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lavanet/lava/utils"
	"github.com/lavanet/lava/utils/sigs"
	pairingtypes "github.com/lavanet/lava/x/pairing/types"
)

func NewRelayFinalizationMetaDataFromReplyMetadataAndRelayRequest(reply ReplyMetadata, req pairingtypes.RelayRequest, consumerAddr sdk.AccAddress) RelayFinalization {
	return RelayFinalization{
		FinalizedBlocksHashes: reply.FinalizedBlocksHashes,
		LatestBlock:           reply.LatestBlock,
		ConsumerAddress:       string(consumerAddr.Bytes()),
		RelaySession:          req.RelaySession,
		SigBlocks:             reply.SigBlocks,
	}
}

func NewRelayFinalizationMetaDataFromRelaySessionAndRelayReply(relaySession *pairingtypes.RelaySession, relayReply *pairingtypes.RelayReply, consumerAddr sdk.AccAddress) RelayFinalization {
	return RelayFinalization{
		FinalizedBlocksHashes: relayReply.FinalizedBlocksHashes,
		LatestBlock:           relayReply.LatestBlock,
		ConsumerAddress:       string(consumerAddr.Bytes()),
		RelaySession:          relaySession,
		SigBlocks:             relayReply.SigBlocks,
	}
}

func (rfm RelayFinalization) GetSignature() []byte {
	return rfm.SigBlocks
}

func (rf RelayFinalization) DataToSign() []byte {
	if rf.RelaySession == nil {
		return nil
	}

	relaySessionHash := tendermintcrypto.Sha256(rf.RelaySession.CalculateHashForFinalization())
	latestBlockBytes := sigs.EncodeUint64(uint64(rf.LatestBlock))
	msgParts := [][]byte{
		latestBlockBytes,
		rf.FinalizedBlocksHashes,
		[]byte(rf.ConsumerAddress),
		relaySessionHash,
	}
	return sigs.Join(msgParts)
}

func (rfm RelayFinalization) HashRounds() int {
	return 1
}

func (rf RelayFinalization) Stringify() string {
	consumerAddr := sdk.AccAddress{}
	consumerAddr.Unmarshal([]byte(rf.ConsumerAddress))
	consumerAddrStr := consumerAddr.String()

	finalizedBlocks := map[int64]string{}
	json.Unmarshal(rf.FinalizedBlocksHashes, &finalizedBlocks)

	return fmt.Sprintf("{latestBlock: %v consumerAddress: %v finalizedBlocksHashes: %v relaySessionHash: %v}",
		utils.StrValue(rf.LatestBlock),
		consumerAddrStr,
		utils.StrValue(finalizedBlocks),
		utils.StrValue(rf.RelaySession.CalculateHashForFinalization()),
	)
}
