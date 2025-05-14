package oauth

import (
	"backend/internal/config"
	"backend/internal/models"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type FatSecretAuthClient struct {
	consumerKey    string
	consumerSecret string
	callbackURL    string
}

func NewFatSecretAuthClient(envs *config.Envs) *FatSecretAuthClient {
	return &FatSecretAuthClient{
		consumerKey:    envs.FatsecretConsumerKey,
		consumerSecret: envs.FatsecretConsumerSecret,
		callbackURL:    envs.FatsecretCallbackURL,
	}
}

func (c *FatSecretAuthClient) GetRequestToken() (string, string, error) {
	params := map[string]string{
		"oauth_consumer_key":     c.consumerKey,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_nonce":            generateNonce(),
		"oauth_version":          "1.0",
		"oauth_callback":         c.callbackURL, //"https://google.com", //c.callbackURL,
	}

	baseString := buildBaseString(
		"POST",
		"https://authentication.fatsecret.com/oauth/request_token",
		params,
	)

	signature := signRequest(baseString, c.consumerSecret, "")
	params["oauth_signature"] = signature

	resp, err := sendOAuthRequest(
		"POST",
		"https://authentication.fatsecret.com/oauth/request_token",
		params,
	)
	if err != nil {
		return "", "", err
	}

	values, err := parseResponse(resp)
	if err != nil {
		return "", "", err
	}

	return values.Get("oauth_token"), values.Get("oauth_token_secret"), nil
}

func (c *FatSecretAuthClient) GetAuthorizationURL(token string) (string, error) {
	authURL := fmt.Sprintf(
		"https://authentication.fatsecret.com/oauth/authorize?oauth_token=%s",
		url.QueryEscape(token),
	)
	return authURL, nil
}

func (c *FatSecretAuthClient) GetAccessToken(requestToken, requestSecret, verifier string) (string, string, error) {
	params := map[string]string{
		"oauth_consumer_key":     c.consumerKey,
		"oauth_token":            requestToken,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_nonce":            generateNonce(),
		"oauth_version":          "1.0",
		"oauth_verifier":         verifier,
	}

	baseString := buildBaseString(
		"POST",
		"https://authentication.fatsecret.com/oauth/access_token",
		params,
	)

	signature := signRequest(baseString, c.consumerSecret, requestSecret)
	params["oauth_signature"] = signature

	resp, err := sendOAuthRequest(
		"POST",
		"https://authentication.fatsecret.com/oauth/access_token",
		params,
	)
	if err != nil {
		return "", "", err
	}

	values, err := parseResponse(resp)
	if err != nil {
		return "", "", err
	}

	return values.Get("oauth_token"), values.Get("oauth_token_secret"), nil
}

func (c *FatSecretAuthClient) GetFoodEntries(ctx context.Context, accessToken, accessSecret string, date time.Time) ([]models.NutritionResponse, error) {
	unixDays := strconv.Itoa(int(date.Unix() / (60 * 60 * 24)))

	baseURL := "https://platform.fatsecret.com/rest/food-entries/v2"

	queryParams := url.Values{}
	queryParams.Add("format", "json")
	queryParams.Add("date", unixDays)
	reqURL := baseURL + "?" + queryParams.Encode()

	oauthParams := map[string]string{
		"format":                 "json",
		"date":                   unixDays,
		"oauth_consumer_key":     c.consumerKey,
		"oauth_token":            accessToken,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_nonce":            generateNonce(),
		"oauth_version":          "1.0",
	}

	baseString := buildBaseString("GET", baseURL, oauthParams)
	signature := signRequest(baseString, c.consumerSecret, accessSecret)
	oauthParams["oauth_signature"] = signature

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var authParts []string
	for k, v := range oauthParams {
		authParts = append(authParts, fmt.Sprintf(`%s="%s"`, k, encode(v)))
	}
	sort.Strings(authParts)
	authHeader := "OAuth " + strings.Join(authParts, ",")
	req.Header.Set("Authorization", authHeader)

	// log.Printf("Final URL: %s", reqURL)
	// log.Printf("Auth Header: %s", authHeader)
	// log.Printf("Base String: %s", baseString)
	// log.Printf("Signature: %s", signature)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Status code error:", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		FoodEntries struct {
			Entry []struct {
				FoodName string `json:"food_entry_name"`
				Calories string `json:"calories"`
				Protein  string `json:"protein"`
				Fat      string `json:"fat"`
				Carbs    string `json:"carbohydrate"`
			} `json:"food_entry"`
		} `json:"food_entries"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Decode json error:", err)
		return nil, err
	}

	if result.Error != nil {
		log.Printf("API error %d: %s\n", result.Error.Code, result.Error.Message)
		return nil, fmt.Errorf("API error %d: %s", result.Error.Code, result.Error.Message)
	}

	entries := make([]models.NutritionResponse, len(result.FoodEntries.Entry))
	for i, item := range result.FoodEntries.Entry {
		calories, err := strconv.ParseFloat(item.Calories, 64)
		if err != nil {
			log.Printf("Invalid calories value '%s': %v", item.Calories, err)
			continue
		}

		protein, err := strconv.ParseFloat(item.Protein, 64)
		if err != nil {
			log.Printf("Invalid protein value '%s': %v", item.Protein, err)
			continue
		}

		fat, err := strconv.ParseFloat(item.Fat, 64)
		if err != nil {
			log.Printf("Invalid fat value '%s': %v", item.Fat, err)
			continue
		}

		carbs, err := strconv.ParseFloat(item.Carbs, 64)
		if err != nil {
			log.Printf("Invalid carbs value '%s': %v", item.Carbs, err)
			continue
		}

		entries[i] = models.NutritionResponse{
			FoodName: item.FoodName,
			Calories: calories,
			Protein:  protein,
			Fat:      fat,
			Carbs:    carbs,
		}
	}

	return entries, nil
}

func generateNonce() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func buildBaseString(method, urlStr string, params map[string]string) string {
	// Collect and sort parameters
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build parameter string
	var paramParts []string
	for _, k := range keys {
		paramParts = append(paramParts, fmt.Sprintf("%s=%s", encode(k), encode(params[k])))
	}
	paramString := strings.Join(paramParts, "&")

	// Build base string
	return fmt.Sprintf("%s&%s&%s",
		encode(method),
		encode(urlStr),
		encode(paramString),
	)
}

func signRequest(baseString, consumerSecret, tokenSecret string) string {
	key := encode(consumerSecret) + "&" + encode(tokenSecret)
	hash := hmac.New(sha1.New, []byte(key))
	hash.Write([]byte(baseString))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func sendOAuthRequest(method, urlStr string, params map[string]string) (*http.Response, error) {
	// Build form data
	form := url.Values{}
	for k, v := range params {
		form.Add(k, v)
	}

	req, err := http.NewRequest(method, urlStr, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	var authParts []string
	for k, v := range params {
		if strings.HasPrefix(k, "oauth_") {
			authParts = append(authParts, fmt.Sprintf(`%s="%s"`, k, encode(v)))
		}
	}
	sort.Strings(authParts)
	authHeader := "OAuth " + strings.Join(authParts, ",")

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	fmt.Println(req.Header.Get("Authorization"))

	return http.DefaultClient.Do(req)
}

func parseResponse(resp *http.Response) (url.Values, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("oauth error: %s", values.Get("oauth_problem"))
	}

	return values, nil
}

func encode(s string) string {
	return strings.ReplaceAll(url.QueryEscape(s), "+", "%20")
}
