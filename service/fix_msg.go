package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/filecoin-project/venus/pkg/crypto"
	"github.com/filecoin-project/venus/venus-shared/types"
)

func (ms *MessageService) FixMsg(ctx context.Context, id string, send bool) error {
	msg, err := ms.GetMessageByUid(ctx, id)
	if err != nil {
		return err
	}
	nonce := msg.Nonce
	estimatedMsg, err := ms.nodeClient.GasEstimateMessageGas(ctx, &msg.Message, &types.MessageSendSpec{}, types.EmptyTSK)
	if err != nil {
		return err
	}
	estimatedMsg.Nonce = nonce
	data, err := estimatedMsg.ToStorageBlock()
	if err != nil {
		return err
	}
	unsignedCid := estimatedMsg.Cid()
	sigI, err := handleTimeout(ms.walletClient.WalletSign, ctx, []interface{}{msg.WalletName, estimatedMsg.From, unsignedCid.Bytes(), types.MsgMeta{
		Type:  types.MTChainMsg,
		Extra: data.RawData(),
	}})

	if err != nil {
		return err
	}
	sig := sigI.(*crypto.Signature)
	result := types.SignedMessage{*estimatedMsg, *sig}

	dataStr, err := json.Marshal(result)
	if err != nil {
		return err
	}
	fmt.Println(string(dataStr))

	if send {
		newCid, err := ms.nodeClient.MpoolPush(ctx, &result)
		if err != nil {
			return err
		}
		fmt.Println("new cid", newCid)
	}
	return nil
}
