package middleware

import (
	"context"
	"mq/common/response"
	"mq/common/util"
	"mq/common/xerr"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginMiddleware struct {
	secret string
}

func NewLoginMiddleware(secret string) *LoginMiddleware {
	return &LoginMiddleware{secret: secret}
}

func (m *LoginMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头中获取 Authorization 字段
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logx.Info("Authorization header is missing")
			response.Response(r, w, nil, xerr.NewErrCode(xerr.InvalidToken))
			return
		}

		// 解析 Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			logx.Info("Bearer token is missing")
			response.Response(r, w, nil, xerr.NewErrCode(xerr.InvalidToken))
			return
		}
		// 解析 JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 确保使用的是预期的签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(m.secret), nil
		})

		if err != nil || !token.Valid {
			logx.Errorf("Invalid token: %v", err)
			response.Response(r, w, nil, xerr.NewErrCode(xerr.InvalidToken))
			return
		}
		// 检查token是否过期
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if exp, ok := claims["exp"].(float64); ok {
				if time.Unix(int64(exp), 0).Before(time.Now()) {
					logx.Error("Token has expired")
					response.Response(r, w, nil, xerr.NewErrCode(xerr.InvalidToken))
					return
				}
			}
		}
		// 从 token 中提取用户信息
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			logx.Error("Invalid token claims")
			response.Response(r, w, nil, xerr.NewErrCode(xerr.InvalidToken))
			return
		}
		// 假设 token 中包含用户 ID
		userId, ok := claims["user_id"]
		if !ok {
			logx.Error("User ID not found in token")
			response.Response(r, w, nil, xerr.NewErrCode(xerr.InvalidToken))
			return
		}

		userIdString := userId.(string)
		logx.Infof("User ID :%s", userIdString)
		// 将 user_id 添加到请求上下文中
		ctx := context.WithValue(r.Context(), util.UserIDKey, userIdString)

		username, ok := claims["username"]
		if !ok {
			logx.Error("Username not found in token")
			response.Response(r, w, nil, xerr.NewErrCode(xerr.InvalidToken))
			return
		}
		usernameString := username.(string)
		logx.Infof("Username :%s", usernameString)
		// 将 username 添加到请求上下文中
		ctx = context.WithValue(ctx, util.UsernameKey, usernameString)
		next(w, r.WithContext(ctx))
	}
}
