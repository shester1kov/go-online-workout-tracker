package oauth

import "backend/internal/config"

type Oauth struct {
	FatSecretAuthClient *FatSecretAuthClient
}

func InitOauth(envs *config.Envs) *Oauth {
	return &Oauth{
		FatSecretAuthClient: NewFatSecretAuthClient(envs),
	}
}
