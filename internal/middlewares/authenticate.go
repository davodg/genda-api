package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/genda/genda-api/internal/app"
	"github.com/genda/genda-api/pkg/config"
	"github.com/julienschmidt/httprouter"

	"github.com/golang-jwt/jwt/v5"
)

func Authenticate(roles []string) app.Middleware {
	currentMw := func(current app.Handler) app.Handler {

		handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
			tokenString := r.Header.Get("Authorization")
			tokenString = strings.Replace(tokenString, "Bearer ", "", -1)

			ok := validateToken(tokenString)
			if !ok {
				data := map[string]string{
					"message": "Not Authorized",
				}

				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(data)
				return nil
			}

			ok = validateJwtScope(tokenString, roles)
			if !ok {
				data := map[string]string{
					"message": "Forbidden",
				}

				w.WriteHeader(http.StatusForbidden)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(data)
				return nil
			}

			return current(ctx, w, r, params)
		}

		return handler
	}

	return currentMw
}

func validateToken(tokenString string) bool {
	if tokenString != "" {
		// if hasTokenAtRedis(rdb, tokenString) {
		// 	return true
		// }
		return validateTokenAtKeycloak(tokenString)
	}

	return false
}

func validateTokenAtKeycloak(tokenString string) bool {
	conf := config.New()
	requestUrl := conf.KeycloakHost + "/protocol/openid-connect/userinfo"
	req, _ := http.NewRequest(http.MethodGet, requestUrl, nil)
	req.Close = true
	req.Header.Add("Authorization", "Bearer "+tokenString)
	req.Header.Add("User-Agent", "gendaApi/1.0.0")

	client := &http.Client{}
	res, err := client.Do(req)

	log.Println(err, res.StatusCode)

	if err == nil && res.StatusCode == http.StatusOK {
		//setTokenAtRedis(rdb, tokenString)
		return true
	}

	return false
}

// func hasTokenAtRedis(rdb *redis.Client, tokenString string) bool {
// 	var ctx = context.Background()
// 	tokenCache, _ := rdb.Get(ctx, tokenString).Result()
// 	return tokenCache != ""
// }

// func setTokenAtRedis(rdb *redis.Client, tokenString string) bool {
// 	var ctx = context.Background()
// 	tokenCache, _ := rdb.Set(ctx, tokenString, tokenString, 10*time.Minute).Result()
// 	return tokenCache != ""
// }

func validateJwtScope(tokenString string, roles []string) bool {
	if len(roles) == 0 {
		return true
	}

	type MyCustomClaims struct {
		RealmAccess struct {
			Roles []string `json:"roles"`
		} `json:"realm_access"`
		jwt.RegisteredClaims
	}

	token, _ := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	claims, _ := token.Claims.(*MyCustomClaims)

	for _, role := range roles {

		ok := slices.Contains(claims.RealmAccess.Roles, role)
		if ok {
			return true
		}
	}

	return false
}
