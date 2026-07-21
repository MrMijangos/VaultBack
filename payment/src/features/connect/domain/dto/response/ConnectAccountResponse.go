package response

import "vault-payment/src/features/connect/domain/entities"

type ConnectAccountResponse struct {
	ChargesEnabled bool `json:"charges_enabled"`
}

func ConnectAccountFromEntity(a *entities.ConnectedAccount) ConnectAccountResponse {
	return ConnectAccountResponse{ChargesEnabled: a.ChargesEnabled}
}
