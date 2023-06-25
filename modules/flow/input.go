package flow

type CreateFlowAccountRequest struct {
	PublicKey string `json:"publicKey" validate:"required"`
}

type CreateFlowAccountResponse struct {
}
