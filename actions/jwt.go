package actions

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	f "teste/models"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user :)
type User struct {
	ID       string `json:"id"`
	Username string `json:"user"`
	Password string `json:"password"` // don't serialize password field
}

// in-memory users
/* var users = []User{
	User{ID: "1", Username: "henrique.b", Password: encryptPassword("T8QKkp5amuSZwvWnEZBYBqsz")},
} */

// LoginRequest represents a login form.
type LoginRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func getUser(username string) (User, error) {
	user := f.AuthUserByUsername(username)[0]
	if user.V_name.String == username {
		teste := User{strconv.Itoa(int(user.V_id.Int64)), user.V_name.String, encryptPassword(user.V_password.String)}
		return teste, nil
	}

	return User{}, errors.New("User not found")
}

func getUserByID(id string) (User, error) {
	user := f.AuthUserById(id)[0]
	if strconv.Itoa(int(user.V_id.Int64)) == id {
		teste := User{strconv.Itoa(int(user.V_id.Int64)), user.V_name.String, encryptPassword(user.V_password.String)}
		return teste, nil
	}

	return User{}, errors.New("User not found")
}

// UsersLogin perform a login with the given credentials.
func UsersLogin(c buffalo.Context) error {

	var req LoginRequest
	err := c.Bind(&req)

	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	pwd := req.Password
	if len(pwd) == 0 {
		return c.Render(http.StatusBadRequest, r.String("Invalid password"))
	}

	username := req.User
	if len(username) == 0 {
		return c.Render(http.StatusBadRequest, r.String("Invalid username"))
	}

	u, err := getUser(username)

	if err != nil {
		return c.Render(http.StatusBadRequest, r.String("Login failed"))
	}

	pwdCompare := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
	if pwdCompare != nil {
		return c.Render(http.StatusBadRequest, r.String("Login failed"))
	}
	// ExpiresAt: time.Now().Add(oneWeek()).Unix(),
	exp := time.Now().Add(5 * time.Minute)

	claims := jwt.StandardClaims{
		ExpiresAt: exp.Unix(),
		Issuer:    fmt.Sprintf("%s.TesteAPI", envy.Get("GO_ENV", "development")),
		Id:        string(u.ID),
	}

	signingKey, err := os.ReadFile(envy.Get("JWT_KEY_PATH", ""))

	if err != nil {
		return fmt.Errorf("could not open jwt key, %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		return fmt.Errorf("could not sign token, %v", err)
	}

	c.Cookies().SetWithExpirationTime("auth", tokenString, exp)

	return c.Render(200, r.JSON(map[string]string{"token": tokenString}))
}

func encryptPassword(p string) string {
	pwd, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(p)), 8)

	if err != nil {
		panic("could not encrypt password")
	}

	return string(pwd)
}

func RestrictedHandlerMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")

		if len(tokenString) == 0 {
			return c.Render(http.StatusForbidden, r.String("No token set in headers"))
		}

		if !strings.Contains(tokenString, "Bearer") {
			return c.Render(http.StatusForbidden, r.String("No Bearer set in Authorization"))
		}

		tokenString = strings.Split(tokenString, "Bearer ")[1]

		// parsing token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// key
			mySignedKey, err := os.ReadFile(envy.Get("JWT_KEY_PATH", ""))

			if err != nil {
				return nil, fmt.Errorf("could not open jwt key, %v", err)
			}

			return mySignedKey, nil
		})

		if err != nil {
			return c.Render(http.StatusForbidden, r.String("Could not parse the token, %v", err))
		}

		// getting claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			logrus.Errorf("claims: %v", claims)

			// retrieving user from db
			u, err := getUserByID(claims["jti"].(string))

			if err != nil {
				return c.Render(http.StatusForbidden, r.String("Could not identify the user"))
			}

			c.Set("user", u)

		} else {
			return c.Render(http.StatusForbidden, r.String("Failed to validate token: %v", claims))
		}

		return next(c)
	}
}
