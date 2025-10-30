package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	gendaConfig "github.com/genda/genda-api/pkg/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/genda/genda-api/internal/app"
	"github.com/julienschmidt/httprouter"
)

type UserTokenSecrets struct {
	KeycloakToken        string `json:"keycloak_token"`
	KeycloakRefreshToken string `json:"keycloak_refresh_token"`
	CalendarToken        string `json:"calendar_token"`
	CalendarRefreshToken string `json:"calendar_refresh_token"`
}

func RefreshKeycloakToken(ctx context.Context, userID string) error {
	conf := gendaConfig.New()
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(conf.AwsAccessKeyId, conf.AwsSecretAccessKey, ""),
		),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return fmt.Errorf("aws config: %w", err)
	}
	sm := secretsmanager.NewFromConfig(awsCfg)
	secretName := "calender/tokens/" + userID

	getOut, err := sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: &secretName})
	if err != nil {
		return fmt.Errorf("get secret: %w", err)
	}

	var fullSecrets UserTokenSecrets
	if err := json.Unmarshal([]byte(*getOut.SecretString), &fullSecrets); err != nil {
		return fmt.Errorf("unmarshal full secret: %w", err)
	}

	data := url.Values{}
	data.Set("client_id", conf.KeycloakClient)
	data.Set("client_secret", conf.ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", fullSecrets.KeycloakRefreshToken)

	tokenURL := fmt.Sprintf("%s/protocol/openid-connect/token", conf.KeycloakHost)
	req, _ := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("keycloak refresh request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("keycloak returned %d: %s", resp.StatusCode, string(body))
	}

	var tok struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Unmarshal(body, &tok); err != nil {
		return fmt.Errorf("unmarshal keycloak resp: %w", err)
	}

	fullSecrets.KeycloakToken = tok.AccessToken
	fullSecrets.KeycloakRefreshToken = tok.RefreshToken

	updatedBytes, err := json.Marshal(fullSecrets)
	if err != nil {
		return fmt.Errorf("marshal updated secret: %w", err)
	}

	_, err = sm.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
		SecretId:     &secretName,
		SecretString: aws.String(string(updatedBytes)),
	})
	if err != nil {
		return fmt.Errorf("put secret: %w", err)
	}

	return nil
}

func InjectSecretToken() app.Middleware {
	return func(next app.Handler) app.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
			raw, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "could not read body", http.StatusBadRequest)
				return nil
			}
			r.Body = io.NopCloser(bytes.NewBuffer(raw))

			var tmp struct {
				UserId *string `json:"user_id"`
			}
			if err := json.Unmarshal(raw, &tmp); err != nil || tmp.UserId == nil {
				http.Error(w, "user_id required", http.StatusBadRequest)
				return nil
			}

			if err := RefreshKeycloakToken(ctx, *tmp.UserId); err != nil {
				http.Error(w, "token refresh failed", http.StatusInternalServerError)
				return nil
			}

			conf := gendaConfig.New()
			awsCfg, _ := config.LoadDefaultConfig(ctx,
				config.WithCredentialsProvider(
					credentials.NewStaticCredentialsProvider(conf.AwsAccessKeyId, conf.AwsSecretAccessKey, ""),
				),
				config.WithRegion("us-east-1"),
			)
			sm := secretsmanager.NewFromConfig(awsCfg)
			getOut, _ := sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: aws.String("calender/tokens/" + *tmp.UserId)})

			var s struct {
				KeycloakToken string `json:"keycloak_token"`
			}
			json.Unmarshal([]byte(*getOut.SecretString), &s)

			r.Header.Set("Authorization", "Bearer "+s.KeycloakToken)

			return next(ctx, w, r, params)
		}
	}
}
