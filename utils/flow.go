package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access"
	"github.com/onflow/flow-go-sdk/access/http"
)

func WaitForSeal(ctx context.Context, c access.Client, id flow.Identifier) (*flow.TransactionResult, error) {
	result, err := c.GetTransactionResult(ctx, id)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Waiting for transaction %s to be sealed...\n", id)

	for result.Status != flow.TransactionStatusSealed {
		time.Sleep(time.Second)
		fmt.Print(".")
		result, err = c.GetTransactionResult(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	fmt.Println()
	fmt.Printf("Transaction %s sealed\n", id)
	return result, nil
}

func GetFlowNetwork(network string) string {
	switch network {
	case "testnet":
		return http.TestnetHost
	case "mainnet":
		return http.MainnetHost
	case "local":
		return http.EmulatorHost
	default:
		return http.EmulatorHost
	}
}
