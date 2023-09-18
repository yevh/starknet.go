package rpc

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type BroadcastedInvokeTransaction interface{}

// AddInvokeTransaction estimates the fee for a given StarkNet transaction.
func (provider *Provider) AddInvokeTransaction(ctx context.Context, broadcastedInvoke BroadcastedInvokeTransaction) (*AddInvokeTransactionResponse, error) {

	var output AddInvokeTransactionResponse
	if err := do(ctx, provider.c, "starknet_addInvokeTransaction", &output, broadcastedInvoke); err != nil {
		return nil, err
	}
	return &output, nil
}

// AddDeclareTransaction submits a new class declaration transaction.
func (provider *Provider) AddDeclareTransaction(ctx context.Context, declareTransaction BroadcastedDeclareTransaction) (*AddDeclareTransactionResponse, error) {
	var result AddDeclareTransactionResponse
	if err := do(ctx, provider.c, "starknet_addDeclareTransaction", &result, declareTransaction); err != nil {
		if strings.Contains(err.Error(), "Invalid contract class") {
			return nil, ErrInvalidContractClass
		}
		return nil, err
	}
	return &result, nil
}

// AddDeployTransaction allows to declare a class and instantiate the
// associated contract in one command. This function will be deprecated and
// replaced by AddDeclareTransaction to declare a class, followed by
// AddInvokeTransaction to instantiate the contract. For now, it remains the only
// way to deploy an account without being charged for it.
func (provider *Provider) AddDeployTransaction(ctx context.Context, deployTransaction BroadcastedDeployTxn) (*AddDeployTransactionResponse, error) {
	var result AddDeployTransactionResponse
	return &result, errors.New("AddDeployTransaction was removed, UDC should be used instead")
}

func (provider *Provider) AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction BroadcastedDeployAccountTransaction) (*AddDeployAccountTransactionResponse, error) {
	fmt.Println("++++++++++++", deployAccountTransaction)
	var result AddDeployAccountTransactionResponse
	if err := do(ctx, provider.c, "starknet_addDeployAccountTransaction", &result, deployAccountTransaction); err != nil {
		fmt.Println("++++++++++++")
		if strings.Contains(err.Error(), "Class hash not found") {
			return nil, ErrClassHashNotFound
		}
		return nil, err
	}
	fmt.Println("++++++++++++-------------")
	return &result, nil
}
