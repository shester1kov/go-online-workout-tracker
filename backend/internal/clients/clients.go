package clients

import "backend/internal/config"

type Clients struct {
	NutritionixClient *NutritionixClient
}

func InitClients(envs *config.Envs) *Clients {
	return &Clients{
		NutritionixClient: NewNutritionixClient(envs),
	}
}
