package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"bigmarket/internal/model"
	"bigmarket/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const adminUsernameKey = "adminUsername"

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func (a *AppHandlers) adminLogin(c *gin.Context) {
	m := params(c)
	username := pStr(m, "username")
	password := pStr(m, "password")
	if username == "" || password == "" {
		render(c, nil, pkg.NewBiz("0002"))
		return
	}
	var user model.AdminUser
	if err := a.App.Store.DB.Where("username = ?", username).First(&user).Error; err != nil || user.Password != hashPassword(password) {
		render(c, nil, pkg.NewBizInfo("0002", "账号或密码错误"))
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(a.Cfg.JWTSecret))
	if err != nil {
		render(c, nil, err)
		return
	}
	render(c, gin.H{"token": signed, "username": user.Username, "nickname": user.Nickname}, nil)
}

func jwtMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, pkg.FailWith("0003", "需要登录"))
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, pkg.FailWith("0003", "登录已失效"))
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if username, ok := claims["username"].(string); ok {
				c.Set(adminUsernameKey, username)
			}
		}
		c.Next()
	}
}

func (a *AppHandlers) adminMe(c *gin.Context) {
	render(c, gin.H{"username": c.GetString(adminUsernameKey)}, nil)
}
