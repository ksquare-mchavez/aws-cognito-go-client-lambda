package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gofiber/fiber/v2"
)

func computeSecretHash(username, clientID, clientSecret string) string {
	key := []byte(clientSecret)
	message := username + clientID
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

type InitiateAuthRequest struct {
	Username string `json:"username"`
}

type InitiateAuthResponse struct {
	ChallengeName string `json:"challenge_name"`
	Session       string `json:"session"`
}

type RespondChallengeRequest struct {
	Username      string `json:"username"`
	OTP           string `json:"otp"`
	ChallengeName string `json:"challenge_name"`
	Session       string `json:"session"`
}

type RespondChallengeResponse struct {
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func getCognitoClient() (*cognitoidentityprovider.Client, string, string) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	clientID := os.Getenv("COGNITO_CLIENT_ID")
	clientSecret := os.Getenv("COGNITO_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Fatalf("COGNITO_CLIENT_ID and COGNITO_CLIENT_SECRET must be set as environment variables")
	}

	return cognitoidentityprovider.NewFromConfig(cfg), clientID, clientSecret
}

func initiateAuthHandler(c *fiber.Ctx) error {
	req := new(InitiateAuthRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	client, clientID, clientSecret := getCognitoClient()
	ctx := context.Background()
	initResp, err := client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: "CUSTOM_AUTH",
		ClientId: aws.String(clientID),
		AuthParameters: map[string]string{
			"USERNAME":    req.Username,
			"SECRET_HASH": computeSecretHash(req.Username, clientID, clientSecret),
		},
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	resp := InitiateAuthResponse{
		ChallengeName: string(initResp.ChallengeName),
		Session:       *initResp.Session,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func respondChallengeHandler(c *fiber.Ctx) error {
	req := new(RespondChallengeRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	client, clientID, clientSecret := getCognitoClient()
	ctx := context.Background()
	resp, err := client.RespondToAuthChallenge(ctx, &cognitoidentityprovider.RespondToAuthChallengeInput{
		ChallengeName: types.ChallengeNameType(req.ChallengeName),
		ClientId:      aws.String(clientID),
		Session:       &req.Session,
		ChallengeResponses: map[string]string{
			"USERNAME":    req.Username,
			"ANSWER":      req.OTP,
			"SECRET_HASH": computeSecretHash(req.Username, clientID, clientSecret),
		},
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	result := RespondChallengeResponse{
		AccessToken:  *resp.AuthenticationResult.AccessToken,
		IdToken:      *resp.AuthenticationResult.IdToken,
		RefreshToken: *resp.AuthenticationResult.RefreshToken,
		TokenType:    *resp.AuthenticationResult.TokenType,
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func main() {
	app := fiber.New()
	app.Post("/initiate-auth", initiateAuthHandler)
	app.Post("/respond-challenge", respondChallengeHandler)
	log.Println("Fiber server started on :8080")
	log.Fatal(app.Listen(":8080"))
}
