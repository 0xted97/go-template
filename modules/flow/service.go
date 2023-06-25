package flow

import (
	"allaccessone/blockchains-support/utils"
	"context"
	"fmt"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/http"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/templates"
)

type Service interface {
	CreateFlowAccount() (bool, string)
}

type FlowService struct {
	ctx        context.Context
	flowClient *http.Client
	// Creator Account
	Signer      crypto.Signer
	AccountKeys []*flow.AccountKey
	AccountKey  *flow.AccountKey
	Account     flow.Address
}

func NewFlowService() *FlowService {
	ctx := context.Background()

	flowClient, err := http.NewClient(utils.GetFlowNetwork(utils.GodotEnv("FLOW_NETWORK")))
	utils.Handle(err)

	privateKey, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, utils.GodotEnv("FLOW_PRIVATE_KEY"))
	utils.Handle(err)

	addr := flow.HexToAddress(utils.GodotEnv("FLOW_ACCOUNT"))
	acc, err := flowClient.GetAccount(context.Background(), addr)
	utils.Handle(err)
	accountKey := acc.Keys[0]

	signer, err := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)
	utils.Handle(err)
	return &FlowService{
		ctx:         ctx,
		flowClient:  flowClient,
		Account:     addr,
		AccountKeys: acc.Keys,
		AccountKey:  accountKey,
		Signer:      signer,
	}
}

func (s *FlowService) UpdateFlowAccount() {
	privateKey, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, utils.GodotEnv("FLOW_PRIVATE_KEY"))
	utils.Handle(err)

	addr := flow.HexToAddress(utils.GodotEnv("FLOW_ACCOUNT"))
	acc, err := s.flowClient.GetAccount(context.Background(), addr)
	utils.Handle(err)
	accountKey := acc.Keys[0]

	signer, err := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)
	utils.Handle(err)

	s.Account = addr
	s.AccountKeys = acc.Keys
	s.AccountKey = accountKey
	s.Signer = signer
}

func (s *FlowService) CreateFlowAccount(input CreateFlowAccountRequest) (flow.Address, error) {
	s.UpdateFlowAccount()

	// seed := []byte("trim action pull regular similar require make weasel biology another banana zebra")
	// privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }
	// publicKey := privateKey.PublicKey()

	publicKey, err := crypto.DecodePublicKeyHex(crypto.ECDSA_secp256k1, input.PublicKey)
	if err != nil {
		return flow.Address{}, err
	}
	accountKey := flow.NewAccountKey().
		SetPublicKey(publicKey).
		SetHashAlgo(crypto.SHA3_256).             // pair this key with the SHA3_256 hashing algorithm
		SetWeight(flow.AccountKeyWeightThreshold) // give this key full signing weight

	tx, _ := templates.CreateAccount([]*flow.AccountKey{accountKey}, nil, s.Account)
	tx.SetProposalKey(
		s.Account,
		s.AccountKey.Index,
		s.AccountKey.SequenceNumber,
	)
	block, err := s.flowClient.GetLatestBlock(context.Background(), true)
	if err != nil {
		return flow.Address{}, err
	}

	tx.SetReferenceBlockID(block.ID)
	tx.SetPayer(s.Account)
	err = tx.SignEnvelope(s.Account, s.AccountKey.Index, s.Signer)
	if err != nil {
		return flow.Address{}, err
	}
	err = s.flowClient.SendTransaction(s.ctx, *tx)
	if err != nil {
		return flow.Address{}, err
	}
	accountCreationTxRes, err := utils.WaitForSeal(s.ctx, s.flowClient, tx.ID())
	if err != nil {
		return flow.Address{}, err
	}
	var myAddress flow.Address

	for _, event := range accountCreationTxRes.Events {
		if event.Type == flow.EventAccountCreated {
			accountCreatedEvent := flow.AccountCreatedEvent(event)
			myAddress = accountCreatedEvent.Address()
		}
	}

	fmt.Println("Account created with address:", myAddress.Hex())

	return myAddress, nil
}
