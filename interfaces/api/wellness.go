package api

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"

	"mq/infrastructure/svc"
	"mq/interfaces/api/handler/appointment"
	"mq/interfaces/api/handler/health"
	"mq/interfaces/api/handler/operation"
	"mq/interfaces/api/handler/patient"
	"mq/interfaces/api/handler/statistics"
	"mq/interfaces/api/handler/tool"
	"mq/interfaces/api/handler/user"
)

// RegisterHandlers 注册HTTP处理器
func RegisterHandlers(server *rest.Server, svc *svc.ServiceContext) {

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/v1/user/login-by-message",
					Handler: user.LoginByMessageHandler(svc),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/user/login-by-token",
					Handler: user.LoginByTokenHandler(svc),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/user/login-by-apple",
					Handler: user.LoginByAppleHandler(svc),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/user/send-message",
					Handler: user.SendMessageHandler(svc),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/user/backend-user",
					Handler: user.BackendUserHandler(svc),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/user/register",
					Handler: user.RegisterHandler(svc),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/user/login",
					Handler: user.LoginHandler(svc),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/user/captcha",
					Handler: user.CaptchaHandler(svc),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/user/send-login-link",
					Handler: user.SendLoginLinkHandler(svc),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/user/email-login",
					Handler: user.EmailLoginHandler(svc),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/user/reset-password",
					Handler: user.ResetPasswordHandler(svc),
				},
			}...,
		),
	)

	// 用户
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{svc.Login},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/v1/user/update-user",
					Handler: user.UpdateUserHandler(svc),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/user/delete-account",
					Handler: user.DeleteAccountHandler(svc),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/user/get-user-info",
					Handler: user.GetUserInfoHandler(svc),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/user/update-password",
					Handler: user.UpdatePasswordHandler(svc),
				},
			}...,
		),
	)

	// 患者管理
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{svc.Login},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/v1/patient/save",
					Handler: patient.SaveHandler(svc),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/patient/get-by-id",
					Handler: patient.GetByIDHandler(svc),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/patient/delete",
					Handler: patient.DeleteHandler(svc),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/patient/list",
					Handler: patient.GetPatientListHandler(svc),
				},
			}...,
		),
	)
	// 健康记录
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{svc.Login},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/v1/health/save",
					Handler: health.SaveHandler(svc),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/health/get-by-id",
					Handler: health.GetByIDHandler(svc),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/health/delete",
					Handler: health.DeleteHandler(svc),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/health/list",
					Handler: health.GetListHandler(svc),
				},
			}...,
		),
	)

	// 预约管理
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{svc.Login},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/v1/appointment/save",
					Handler: appointment.SaveHandler(svc),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/appointment/get-by-id",
					Handler: appointment.GetByIDHandler(svc),
				},
				{
					Method:  http.MethodPost,
					Path:    "/v1/appointment/delete",
					Handler: appointment.DeleteHandler(svc),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/appointment/list",
					Handler: appointment.GetListHandler(svc),
				},
			}...,
		),
	)
	// 操作记录
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{svc.Login},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/v1/operation/save",
					Handler: operation.CreateHandler(svc),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/operation/list",
					Handler: operation.GetListHandler(svc),
				},
			}...,
		),
	)
	// 统计
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{svc.Login},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/v1/statistics/now",
					Handler: statistics.GetNowStatisticsHandler(svc),
				},
			}...,
		),
	)

	// 工具
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/v1/tool/upload-token",
				Handler: tool.UploadTokenHandler(svc),
			},
		},
	)
}
