package main

import (
	"context"
	"fmt"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access"
	"github.com/onflow/flow-go-sdk/access/http"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/templates"
)

// https://github.com/restuwahyu13/go-rest-api/blob/main/utils/response.go

func Handle(err error) {
	if err != nil {
		fmt.Println("err:", err.Error())
		panic(err)
	}
}

func ServiceAccount(flowClient access.Client) (flow.Address, *flow.AccountKey, crypto.Signer) {
	privateKey, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, "f93e4891a468030319a5e86eaa1204d1dbe7a534632a56720c11d63cd4c0739e")
	Handle(err)

	addr := flow.HexToAddress("0xb941442fdd844a30")
	acc, err := flowClient.GetAccount(context.Background(), addr)
	Handle(err)
	fmt.Printf("acc: %v\n", acc.Keys)
	accountKey := acc.Keys[0]
	signer, err := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)
	Handle(err)
	return addr, accountKey, signer
}

func No_main() {
	ctx := context.Background()

	seed := []byte("trim action pull regular similar require make weasel biology another banana zebra")
	privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	flowClient, err := http.NewClient(http.TestnetHost)
	if err != nil {
		panic("failed to connect to emulator")
	}
	block, err := flowClient.GetLatestBlock(context.Background(), true)
	Handle(err)
	fmt.Printf("block: %v\n", block.ID)
	// get the public key
	publicKey := privateKey.PublicKey()
	accountKey := flow.NewAccountKey().
		SetPublicKey(publicKey).
		SetHashAlgo(crypto.SHA3_256).             // pair this key with the SHA3_256 hashing algorithm
		SetWeight(flow.AccountKeyWeightThreshold) // give this key full signing weight

	serviceAcctAddr, serviceAcctKey, serviceSigner := ServiceAccount(flowClient)
	tx, _ := templates.CreateAccount([]*flow.AccountKey{accountKey}, nil, serviceAcctAddr)
	// // connect to an emulator running locally
	tx.SetProposalKey(
		serviceAcctAddr,
		serviceAcctKey.Index,
		serviceAcctKey.SequenceNumber,
	)
	tx.SetReferenceBlockID(block.ID)
	tx.SetPayer(serviceAcctAddr)
	fmt.Printf("tx.Authorizers: %v\n", tx.Authorizers)
	err = tx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	Handle(err)

	err = flowClient.SendTransaction(ctx, *tx)
	Handle(err)

}
