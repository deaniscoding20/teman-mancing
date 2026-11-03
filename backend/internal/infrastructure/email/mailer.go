package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Config struct {
	Host        string
	Port        int
	User        string
	Password    string
	UseSSL      bool
	SenderEmail string
	SenderName  string
	BaseURL     string
}

type Mailer struct {
	cfg    Config
	logger *zap.Logger
}

func NewMailer(cfg Config, logger *zap.Logger) *Mailer {
	if cfg.Port == 0 {
		cfg.Port = 465
	}
	if cfg.SenderEmail == "" {
		cfg.SenderEmail = "noreply@sumopod.com"
	}
	if cfg.SenderName == "" {
		cfg.SenderName = "Teman Mancing"
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://temanmancing.my.id"
	}
	return &Mailer{
		cfg:    cfg,
		logger: logger,
	}
}

func (m *Mailer) SendMail(toEmail, subject, htmlBody string) error {
	if m.cfg.Host == "" || m.cfg.User == "" {
		m.logger.Warn("SMTP not configured, skipping email dispatch", zap.String("to", toEmail))
		return nil
	}

	fromHeader := fmt.Sprintf("%s <%s>", m.cfg.SenderName, m.cfg.SenderEmail)
	headers := make(map[string]string)
	headers["From"] = fromHeader
	headers["To"] = toEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"
	headers["Date"] = time.Now().Format(time.RFC1123Z)

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)

	if m.cfg.Port == 465 || m.cfg.UseSSL {
		tlsConfig := &tls.Config{
			ServerName: m.cfg.Host,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to dial SMTP SSL: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, m.cfg.Host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer client.Quit()

		auth := smtp.PlainAuth("", m.cfg.User, m.cfg.Password, m.cfg.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate SMTP: %w", err)
		}

		if err := client.Mail(m.cfg.SenderEmail); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		if err := client.Rcpt(toEmail); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}

		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to initiate data stream: %w", err)
		}
		defer w.Close()

		if _, err := w.Write([]byte(msg.String())); err != nil {
			return fmt.Errorf("failed to write email body: %w", err)
		}

		m.logger.Info("Email successfully sent via SSL", zap.String("to", toEmail), zap.String("subject", subject))
		return nil
	}

	// Port 587 STARTTLS fallback
	auth := smtp.PlainAuth("", m.cfg.User, m.cfg.Password, m.cfg.Host)
	err := smtp.SendMail(addr, auth, m.cfg.SenderEmail, []string{toEmail}, []byte(msg.String()))
	if err != nil {
		return fmt.Errorf("failed to send email via STARTTLS: %w", err)
	}

	m.logger.Info("Email successfully sent via STARTTLS", zap.String("to", toEmail), zap.String("subject", subject))
	return nil
}

func (m *Mailer) SendVerificationEmail(toEmail, toName, otpCode, verifyToken string) error {
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", strings.TrimRight(m.cfg.BaseURL, "/"), verifyToken)
	subject := fmt.Sprintf("Verifikasi Email Anda - Kode OTP: %s", otpCode)

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #f8fafc; color: #1e293b; margin: 0; padding: 24px; }
    .card { max-width: 520px; margin: 0 auto; background: #ffffff; border-radius: 20px; padding: 32px; border: 1px solid #e2e8f0; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05); }
    .brand { font-size: 20px; font-weight: 800; color: #0284c7; text-align: center; margin-bottom: 24px; }
    .otp-box { background: #f0fdf4; border: 2px dashed #22c55e; border-radius: 16px; padding: 18px; text-align: center; margin: 24px 0; }
    .otp-code { font-size: 32px; font-weight: 900; letter-spacing: 6px; color: #15803d; font-family: monospace; }
    .btn { display: block; width: 100%%; max-width: 280px; margin: 24px auto; background-color: #0284c7; color: #ffffff !important; text-align: center; padding: 14px 20px; border-radius: 12px; font-weight: 700; text-decoration: none; font-size: 14px; }
    .footer { font-size: 11px; color: #94a3b8; text-align: center; margin-top: 24px; }
  </style>
</head>
<body>
  <div class="card">
    <div class="brand">🎣 Teman Mancing</div>
    <h2 style="font-size: 18px; font-weight: 700; margin-bottom: 8px;">Halo, %s!</h2>
    <p style="font-size: 14px; color: #64748b; line-height: 1.6;">
      Terima kasih telah mendaftar di <strong>Teman Mancing</strong>. Silakan masukkan kode OTP berikut atau klik tombol untuk mengaktifkan akun Anda:
    </p>

    <div class="otp-box">
      <div style="font-size: 12px; color: #166534; font-weight: 600; margin-bottom: 4px;">KODE VERIFIKASI (OTP)</div>
      <div class="otp-code">%s</div>
      <div style="font-size: 11px; color: #16a34a; margin-top: 4px;">Berlaku selama 15 menit</div>
    </div>

    <a href="%s" class="btn">Verifikasi Akun Sekarang</a>

    <p style="font-size: 12px; color: #94a3b8; line-height: 1.5; margin-top: 20px;">
      Jika tombol di atas tidak berfungsi, buka tautan berikut di browser Anda:<br>
      <a href="%s" style="color: #0284c7; word-break: break-all;">%s</a>
    </p>

    <div class="footer">
      Jika Anda tidak merasa mendaftar di Teman Mancing, abaikan email ini.<br>
      &copy; %d Teman Mancing. Hak cipta dilindungi.
    </div>
  </div>
</body>
</html>`, toName, otpCode, verifyURL, verifyURL, verifyURL, time.Now().Year())

	return m.SendMail(toEmail, subject, body)
}

func (m *Mailer) SendPasswordResetEmail(toEmail, toName, otpCode, resetToken string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", strings.TrimRight(m.cfg.BaseURL, "/"), resetToken)
	subject := fmt.Sprintf("Atur Ulang Kata Sandi - Kode: %s", otpCode)

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #f8fafc; color: #1e293b; margin: 0; padding: 24px; }
    .card { max-width: 520px; margin: 0 auto; background: #ffffff; border-radius: 20px; padding: 32px; border: 1px solid #e2e8f0; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05); }
    .brand { font-size: 20px; font-weight: 800; color: #0284c7; text-align: center; margin-bottom: 24px; }
    .otp-box { background: #fef2f2; border: 2px dashed #ef4444; border-radius: 16px; padding: 18px; text-align: center; margin: 24px 0; }
    .otp-code { font-size: 32px; font-weight: 900; letter-spacing: 6px; color: #b91c1c; font-family: monospace; }
    .btn { display: block; width: 100%%; max-width: 280px; margin: 24px auto; background-color: #ef4444; color: #ffffff !important; text-align: center; padding: 14px 20px; border-radius: 12px; font-weight: 700; text-decoration: none; font-size: 14px; }
    .footer { font-size: 11px; color: #94a3b8; text-align: center; margin-top: 24px; }
  </style>
</head>
<body>
  <div class="card">
    <div class="brand">🎣 Teman Mancing</div>
    <h2 style="font-size: 18px; font-weight: 700; margin-bottom: 8px;">Permintaan Atur Ulang Password</h2>
    <p style="font-size: 14px; color: #64748b; line-height: 1.6;">
      Halo <strong>%s</strong>, kami menerima permintaan untuk mengatur ulang kata sandi akun Anda. Gunakan kode verifikasi di bawah atau klik tombol:
    </p>

    <div class="otp-box">
      <div style="font-size: 12px; color: #991b1b; font-weight: 600; margin-bottom: 4px;">KODE RESET PASSWORD</div>
      <div class="otp-code">%s</div>
      <div style="font-size: 11px; color: #ef4444; margin-top: 4px;">Berlaku selama 15 menit</div>
    </div>

    <a href="%s" class="btn">Atur Ulang Password Sekarang</a>

    <p style="font-size: 12px; color: #94a3b8; line-height: 1.5; margin-top: 20px;">
      Jika tombol di atas tidak berfungsi, buka tautan berikut:<br>
      <a href="%s" style="color: #0284c7; word-break: break-all;">%s</a>
    </p>

    <div class="footer">
      Jika Anda tidak meminta perubahan password, amankan akun Anda atau abaikan pesan ini.<br>
      &copy; %d Teman Mancing. Hak cipta dilindungi.
    </div>
  </div>
</body>
</html>`, toName, otpCode, resetURL, resetURL, resetURL, time.Now().Year())

	return m.SendMail(toEmail, subject, body)
}
