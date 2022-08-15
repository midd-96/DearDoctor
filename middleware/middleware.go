package middleware

import (
	"dearDoctor/common/response"
	"dearDoctor/service"
	"dearDoctor/utils"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Middleware interface {
	AuthorizeJwt(http.Handler) http.Handler
}

type middleware struct {
	jwtService service.JWTService
}

func NewMiddlewareAdmin(jwtAdminService service.JWTService) Middleware {
	return &middleware{
		jwtService: jwtAdminService,
	}

}

func NewMiddlewareUser(jwtUserService service.JWTService) Middleware {
	return &middleware{
		jwtService: jwtUserService,
	}

}

func NewMiddlewareDoctors(jwtDoctorService service.JWTService) Middleware {
	return &middleware{
		jwtService: jwtDoctorService,
	}

}

func (c *middleware) AuthorizeJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//getting from header
		autheader := r.Header.Get("Authorization")
		bearerToken := strings.Split(autheader, " ")

		if len(bearerToken) != 2 {
			err := errors.New("request does not contain an access token")
			response := response.ErrorResponse("Failed to autheticate jwt", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			utils.ResponseJSON(w, response)
			return
		}

		authtoken := bearerToken[1]
		ok, claims := c.jwtService.VerifyToken(authtoken)

		if !ok {
			err := errors.New("your token is not valid")
			response := response.ErrorResponse("Error", err.Error(), nil)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			utils.ResponseJSON(w, response)
			return
		}

		user_id := fmt.Sprint(claims.User_Id)
		r.Header.Set("user_id", user_id)
		next.ServeHTTP(w, r)

	})
}
