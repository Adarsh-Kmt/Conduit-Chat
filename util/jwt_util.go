package util

import (
	"net/http"
	"os"
	"time"

	"fmt"
	"log"

	"github.com/Adarsh-Kmt/chatapp/types"

	jwt "github.com/golang-jwt/jwt/v5"
)

func MakeJWTAuthHttpHandlerFunc(f http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Auth")

		// if substr := tokenString[:7]; len(tokenString) == 0 || substr != "Bearer " {

		// 	WriteJSON(w, http.StatusForbidden, map[string]string{"Error": "incorrect auth token."})
		// 	return

		// }
		log.Println(tokenString)
		_, err := ValidateJWTToken(tokenString)

		if err != nil {

			WriteJSON(w, http.StatusForbidden, map[string]string{"Error": err.Error()})
			return
		}

		f(w, r)

	}
}

func GenerateJWTToken(UserObjectId string) (string, *types.APIError) {

	// Parse the RSA private key
	//pkcs1 format
	// pemData, err := os.ReadFile("private_key.pem")
	// if err != nil {
	// 	fmt.Println("Error reading PEM file:", err)
	// 	os.Exit(1)
	// }
	// privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemData)

	secretKey := os.Getenv("JWT_PRIVATE_KEY")
	// if err != nil {
	// 	log.Printf(err.Error())
	// 	return "", &APIError{Error: "error while parsing private key.", ErrorStatus: 500}
	// }

	claims := &jwt.RegisteredClaims{
		Subject:  UserObjectId,
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	log.Println(secretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		log.Println(err.Error())
		return "", &types.APIError{Error: "error while generating jwt token.", ErrorStatus: 500}
	}

	return tokenString, nil
}

func ValidateJWTToken(tokenString string) (*jwt.Token, error) {

	// publicKeyBytes, err := os.ReadFile("public_key.pem")
	// if err != nil {
	// 	return nil, fmt.Errorf("error while reading public key: %v", err)
	// }

	// PublicKeyNew, err := x509.Pa(publicKeyBytes)

	// if err != nil {
	// 	log.Println(err.Error())
	// 	return nil, err
	// }
	// // Parse the PEM-encoded public key
	// block, _ := pem.Decode(publicKeyBytes)
	// if block == nil {
	// 	return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	// }
	// if block.Type != "PUBLIC KEY" {
	// 	return nil, fmt.Errorf("unexpected PEM type: %s", block.Type)
	// }

	// // Parse the public key
	// publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	// if err != nil {
	// 	return nil, fmt.Errorf("error parsing public key: %v", err)
	// }

	// if _, ok := publicKey.(*rsa.PublicKey); !ok {

	// 	return nil, fmt.Errorf("wrong public key.")
	// }
	log.Println("the token is")
	log.Println(tokenString)

	secretKey := os.Getenv("JWT_PRIVATE_KEY")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("wrong signing key.")
		}

		return []byte(secretKey), nil
	})

}

func GetUserObjectIdFromJWT(JWTToken string) (string, error) {

	secretKey := os.Getenv("JWT_PRIVATE_KEY")

	parsedToken, err := jwt.Parse(JWTToken, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("wrong signing key.")
		}

		return []byte(secretKey), nil
	})

	if err != nil {

		return "", fmt.Errorf("error in parsing JWT")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", fmt.Errorf("error in getting claims from JWT.")
	}

	userObjectID, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("error in extracting subject claim from JWT.")

	}

	return userObjectID, nil

}
