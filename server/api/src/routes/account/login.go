package account

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"../../../../lib"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UUID     string `json:"uuid"`
}

func generateRandomSHA256() (string, string) {
	hash := sha256.New()
	generated := lib.GetRandomString(43)
	hash.Write([]byte(generated))
	hashString := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return hashString, generated
}

func checkUserSecret(inputData loginData, db *sqlx.DB) (lib.User, int, string) {
	var u lib.User
	err := db.Get(&u, "SELECT * FROM Users WHERE username = $1", inputData.Username)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(lib.PrettyError("[DB REQUEST - SELECT] Failed to get user data " + err.Error()))
		return lib.User{}, 500, "User data collection failed"
	}
	if err == sql.ErrNoRows || u == (lib.User{}) {
		return lib.User{}, 403, "User or password incorrect"
	}
	// Comparing the password with the hashed password from the body
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputData.Password))
	if err != nil {
		return lib.User{}, 403, "User or password incorrect"
	}
	return u, 0, ""
}

func generateJWT(u lib.User, UUID string, client *redis.Client) (string, int, string) {
	now := time.Now().Local()
	duration := time.Hour * time.Duration(72)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":      "matcha.com",
		"sub":      UUID,
		"userId":   u.ID,
		"username": u.Username,
		"iat":      now.Unix(),
		"exp":      now.Add(duration).Unix(),
	})
	tokenString, err := token.SignedString(lib.JWTSecret)
	if err != nil {
		return "", 500, "JWT creation failed"
	}
	err = lib.RedisSetValue(client, u.Username+"-"+UUID, tokenString, duration)
	if err != nil {
		return "", 500, "Insert key:value in Redis failed"
	}
	return tokenString, 0, ""
}

// Login function corresponds to the API route /v1/account/login
// It allows the handle the user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.Database).(*sqlx.DB)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	redisClient, ok := r.Context().Value(lib.Redis).(*redis.Client)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with redis connection")
		return
	}
	var inputData loginData
	errCode, errContent, err := lib.GetDataBody(r, &inputData)
	if err != nil {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	u, errCode, errContent := checkUserSecret(inputData, db)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	// Create JSON Web Token
	token, errCode, errContent := generateJWT(u, inputData.UUID, redisClient)
	if errCode != 0 || errContent != "" {
		lib.RespondWithErrorHTTP(w, errCode, errContent)
		return
	}
	lib.RespondWithJSON(w, 200, map[string]interface{}{
		"token": token,
	})
}
