package models

type FatSecretAuth struct {
	UserID       int    `json:"user_id"`
	AccessToken  string `json:"access_token"`
	AccessSecret string `json:"access_secret"`
}

type TempAuth struct {
	UserID        int    `json:"user_id"`
	RequestToken  string `json:"request_token"`
	RequestSecret string `json:"request_secret"`
}
