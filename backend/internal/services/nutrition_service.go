package services

import (
	"backend/internal/apperrors"
	"backend/internal/models"
	"backend/internal/oauth"
	"backend/internal/repository"
	"context"
	"log"
	"net/http"
	"time"
)

type NutritionService struct {
	authRepo            *repository.FatSecretAuthRepository
	fatSecretAuthClient *oauth.FatSecretAuthClient
}

func NewNutritionService(
	authRepo *repository.FatSecretAuthRepository,
	fatSecretAuthClient *oauth.FatSecretAuthClient,
) *NutritionService {
	return &NutritionService{
		authRepo:            authRepo,
		fatSecretAuthClient: fatSecretAuthClient,
	}
}

func (s *NutritionService) GetDailyNutrition(ctx context.Context, date time.Time) (*[]models.NutritionEntry, error) {
	userID := ctx.Value("user_id").(int)

	auth, err := s.authRepo.GetAuth(ctx, userID)
	if err != nil {
		return nil, err
	}

	// client := clients.NewFatSecretClient(auth.AccessToken, auth.AccessSecret, s.fatSecretAuthClient)
	fsEntries, err := s.fatSecretAuthClient.GetFoodEntries(ctx, auth.AccessToken, auth.AccessSecret, date)
	if err != nil {
		return nil, err
	}

	var entries []models.NutritionEntry

	for _, fsEntry := range fsEntries {
		entry := models.NutritionEntry{
			UserID:   userID,
			FoodName: fsEntry.FoodName,
			Calories: fsEntry.Calories,
			Protein:  fsEntry.Protein,
			Fat:      fsEntry.Fat,
			Carbs:    fsEntry.Carbs,
			Date:     date,
		}

		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		log.Println("No food entries found for the date:", date.Format("2006-01-02"))
		return nil, &apperrors.AppError{
			Code:    http.StatusNotFound,
			Message: "Нет записей о еде на эту дату",
		}
	}

	return &entries, nil
}

func (s *NutritionService) InitFatSecretAuth(ctx context.Context) (string, error) {
	userID := ctx.Value("user_id").(int)

	requestToken, requestSecret, err := s.fatSecretAuthClient.GetRequestToken()
	if err != nil {
		log.Println("Error get request token:", err)
		return "", err
	}

	log.Println("req token, req secret:", requestToken, requestSecret)
	err = s.authRepo.SaveTempAuth(ctx, &models.TempAuth{
		UserID:        userID,
		RequestToken:  requestToken,
		RequestSecret: requestSecret,
	})
	if err != nil {
		log.Println("Save temp auth error:", err)
		return "", err
	}

	authURL, err := s.fatSecretAuthClient.GetAuthorizationURL(requestToken)
	if err != nil {
		log.Println("Error get auth url")
		return "", err
	}

	log.Println("auth url:", authURL)
	return authURL, nil
}

func (s *NutritionService) CompleteFatSecretAuth(ctx context.Context, token, verifier string) error {
	tempAuth, err := s.authRepo.GetTempAuth(ctx, token)
	if err != nil {
		log.Println("Error get temp auth:", err)
		return err
	}

	accessToken, accessSecret, err := s.fatSecretAuthClient.GetAccessToken(
		tempAuth.RequestToken,
		tempAuth.RequestSecret,
		verifier,
	)
	if err != nil {
		log.Println("Error get access token:", err)
		return err
	}

	return s.authRepo.SaveAuth(ctx, &models.FatSecretAuth{
		UserID:       tempAuth.UserID,
		AccessToken:  accessToken,
		AccessSecret: accessSecret,
	})
}
