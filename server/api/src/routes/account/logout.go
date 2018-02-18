package account

import (
	"net/http"

	"../../../../lib"

	"github.com/go-redis/redis"
)

// Logout is the route '/v1/account/logout' with the method POST.
// Delete in the Redis database the key `Username + "-" + UUID` (using context data)
// If deletion failed
//    -> Return an error - HTTP Code 500 Internal Server Error - JSON Content "Error: Failed to delete token"
// Return HTTP Code 202 Status Accepted
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
	lib.RespondEmptyHTTP(w, http.StatusAccepted)
}
