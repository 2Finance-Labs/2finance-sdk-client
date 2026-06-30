package client_2finance

import (
	"encoding/json"
	"fmt"

	"github.com/2Finance-Labs/2finance-sdk-client/protocol"
	"gitlab.com/2finance/2finance-network/blockchain/types"
	"gitlab.com/2finance/2finance-network/blockchain/virtualmachine"
)

func (c *NetworkClient) SignPreparedTransaction(tx protocol.PreparedTransaction) (protocol.SignedTransaction, error) {
	if c.walletManager == nil {
		return protocol.SignedTransaction{}, fmt.Errorf("wallet manager is required")
	}

	return c.walletManager.SignPreparedTransaction(tx)
}

func (c *NetworkClient) SubmitSignedTransaction(tx protocol.SignedTransaction) (types.ContractOutput, error) {
	networkTx, err := tx.ToNetworkTransaction()
	if err != nil {
		return types.ContractOutput{}, err
	}

	contractOutputBytes, err := c.SendTransaction(
		virtualmachine.REQUEST_METHOD_SEND,
		networkTx,
		c.replyTo,
	)
	if err != nil {
		return types.ContractOutput{}, fmt.Errorf("failed to send transaction: %w", err)
	}

	var contractOutput types.ContractOutput
	if err := json.Unmarshal(contractOutputBytes, &contractOutput); err != nil {
		return types.ContractOutput{}, fmt.Errorf("failed to unmarshal contract output: %w", err)
	}

	return contractOutput, nil
}

func (c *NetworkClient) SignAndSendPreparedTransaction(tx protocol.PreparedTransaction) (types.ContractOutput, error) {
	signed, err := c.SignPreparedTransaction(tx)
	if err != nil {
		return types.ContractOutput{}, err
	}

	return c.SubmitSignedTransaction(signed)
}
