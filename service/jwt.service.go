package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTService interface {
	GenerateToken(user_id int, username string, role string) string
	VerifyToken(token string) (bool, *SignedDetails)
	GetTokenFromString(signedToken string, claims *SignedDetails) (*jwt.Token, error)
	GenerateRefreshToken(token string) (string, error)
}

type SignedDetails struct {
	User_Id  int
	Username string
	Role     string
	jwt.StandardClaims
}

type jwtService struct {
	SecretKey string
}

func NewJWTUserService() JWTService {
	return &jwtService{
		SecretKey: os.Getenv("USER_KEY"),
	}
}

func NewJWTDoctorService() JWTService {
	return &jwtService{
		SecretKey: os.Getenv("DOCTOR_KEY"),
	}
}

func NewJWTAdminService() JWTService {
	return &jwtService{
		SecretKey: os.Getenv("ADMIN_KEY"),
	}
}

func (j *jwtService) GenerateToken(userId int, email, role string) string {
	claims := &SignedDetails{
		userId,
		email,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.SecretKey))

	if err != nil {
		log.Println(err)
	}

	return signedToken
}
func (j *jwtService) GenerateRefreshToken(accessToken string) (string, error) {
	claims := &SignedDetails{}
	j.GetTokenFromString(accessToken, claims)

	if time.Until(time.Unix(claims.ExpiresAt, 0)) > 30*time.Second {
		return "", errors.New("too early to generate refresh token")
	}

	claims.ExpiresAt = time.Now().Local().Add(time.Minute * time.Duration(5)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshToken, err := token.SignedString([]byte(j.SecretKey))

	if err != nil {
		log.Println(err)
	}
	return refreshToken, err

}

func (j *jwtService) VerifyToken(signedToken string) (bool, *SignedDetails) {
	claims := &SignedDetails{}
	token, _ := j.GetTokenFromString(signedToken, claims)

	if token.Valid {
		if e := claims.Valid(); e == nil {
			return true, claims
		}
	}
	return false, claims
}

func (j *jwtService) GetTokenFromString(signedToken string, claims *SignedDetails) (*jwt.Token, error) {

	return jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.SecretKey), nil
	})

}
