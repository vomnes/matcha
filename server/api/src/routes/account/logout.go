package account

import (
	"net/http"

	"../../../../lib"

	"github.com/go-redis/redis"
)

// Logout function corresponds to the API route POST /v1/account/logout
// It allows to handle the user logout
// Delete the user JSON Web Token from Redis database
func Logout(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(lib.Username).(string)
	UUID := r.Context().Value(lib.UUID).(string)
	redisClient, ok := r.Context().Value(lib.Redis).(*redis.Client)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with redis connection")
		return
	}
	if ok = lib.RedisDelValue(redisClient, username+"-"+UUID); !ok {
		lib.RespondWithErrorHTTP(w, 500, "Failed to delete token")
		return
	}
	lib.RespondWithJSON(w, http.StatusAccepted, "OK")
}
