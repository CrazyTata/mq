package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"mq/infrastructure/svc"
	"net/smtp"

	"encoding/base64"

	"github.com/zeromicro/go-zero/core/logx"
)

var _ EmailInterface = &EmailSender{}

type EmailSender struct {
	SvcCtx *svc.ServiceContext
}

func (s *EmailSender) SendLoginLink(ctx context.Context, emailTo, token string) error {
	// 构建登录链接
	loginLink := fmt.Sprintf("%s/sign-in?token=%s", s.SvcCtx.GetConfig().Email.FrontendURL, token)

	// 邮件内容
	subject := "Login Link"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h3>Confirm your signup</h3>
			<p>Follow this link to confirm your user:</p>
			<p><a href="%s">Confirm your mail</a></p>
			<p>If you did not request this email, please ignore it.</p>
			<p>The link is valid for 15 minutes.</p>
		</body>
		</html>
	`, loginLink)

	// 构建邮件头
	headers := make(map[string]string)
	headers["From"] = s.SvcCtx.GetConfig().Email.From
	headers["To"] = emailTo
	headers["Subject"] = "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(subject)) + "?="
	headers["Content-Type"] = "text/html; charset=UTF-8"
	headers["MIME-Version"] = "1.0"

	// 组装邮件内容
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// 发送邮件
	auth := smtp.PlainAuth("", s.SvcCtx.GetConfig().Email.Username, s.SvcCtx.GetConfig().Email.Password, s.SvcCtx.GetConfig().Email.Host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.SvcCtx.GetConfig().Email.Host,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.SvcCtx.GetConfig().Email.Host, s.SvcCtx.GetConfig().Email.Port), tlsConfig)
	if err != nil {
		logx.Errorf("连接邮件服务器失败: %v", err)
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.SvcCtx.GetConfig().Email.Host)
	if err != nil {
		logx.Errorf("创建SMTP客户端失败: %v", err)
		return err
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		logx.Errorf("SMTP认证失败: %v", err)
		return err
	}

	if err = client.Mail(s.SvcCtx.GetConfig().Email.From); err != nil {
		logx.Errorf("设置发件人失败: %v", err)
		return err
	}

	if err = client.Rcpt(emailTo); err != nil {
		logx.Errorf("设置收件人失败: %v", err)
		return err
	}

	w, err := client.Data()
	if err != nil {
		logx.Errorf("准备发送数据失败: %v", err)
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		logx.Errorf("写入邮件内容失败: %v", err)
		return err
	}

	err = w.Close()
	if err != nil {
		logx.Errorf("关闭数据写入失败: %v", err)
		return err
	}

	return nil
}

func (s *EmailSender) SendResetPwdLink(ctx context.Context, emailTo, token string) error {
	// 构建重置密码链接
	resetLink := fmt.Sprintf("%s/login/reset?token=%s", s.SvcCtx.GetConfig().Email.FrontendURL, token)

	// 邮件内容
	subject := "Reset Password"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h3>Hello!</h3>
			<p>Please click the following link to reset your password:</p>
			<p><a href="%s">Click here to reset your password</a></p>
			<p>If you did not request this email, please ignore it.</p>
			<p>The link is valid for 15 minutes.</p>
		</body>
		</html>
	`, resetLink)

	// 构建邮件头
	headers := make(map[string]string)
	headers["From"] = s.SvcCtx.GetConfig().Email.From
	headers["To"] = emailTo
	headers["Subject"] = "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(subject)) + "?="
	headers["Content-Type"] = "text/html; charset=UTF-8"
	headers["MIME-Version"] = "1.0"

	// 组装邮件内容
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// 发送邮件
	auth := smtp.PlainAuth("", s.SvcCtx.GetConfig().Email.Username, s.SvcCtx.GetConfig().Email.Password, s.SvcCtx.GetConfig().Email.Host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.SvcCtx.GetConfig().Email.Host,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.SvcCtx.GetConfig().Email.Host, s.SvcCtx.GetConfig().Email.Port), tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.SvcCtx.GetConfig().Email.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(s.SvcCtx.GetConfig().Email.From); err != nil {
		return err
	}

	if err = client.Rcpt(emailTo); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}
