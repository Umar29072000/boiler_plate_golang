package email

import (
	"boiler_plate_be_golang/app/config"
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"path/filepath"
)

// EmailService handles email sending
type EmailService struct {
	From      string
	Host      string
	Port      string
	Username  string
	Password  string
	AppConfig config.App
}

// NewEmailService creates a new email service
func NewEmailService(emailConfig config.Email, appConfig config.App) *EmailService {
	return &EmailService{
		From:      emailConfig.From,
		Host:      emailConfig.Host,
		Port:      emailConfig.Port,
		Username:  emailConfig.Username,
		Password:  emailConfig.Password,
		AppConfig: appConfig,
	}
}

// EmailData represents data for email templates
type EmailData struct {
	Name            string
	Email           string
	VerificationURL string
	ResetURL        string
	AppName         string
	AppURL          string
}

// SendEmail sends an email using SMTP
func (s *EmailService) SendEmail(to, subject, htmlBody string) error {
	// Setup email headers
	headers := make(map[string]string)
	headers["From"] = s.From
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	// Setup authentication
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)

	// For production SMTP servers with TLS
	if s.AppConfig.Env == "production" {
		// Setup TLS config
		tlsconfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.Host,
		}

		// Connect to server
		conn, err := tls.Dial("tcp", addr, tlsconfig)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer conn.Close()

		// Create SMTP client
		client, err := smtp.NewClient(conn, s.Host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer client.Close()

		// Authenticate
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		// Set sender and recipient
		if err = client.Mail(s.From); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}

		// Send email body
		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to get data writer: %w", err)
		}
		_, err = w.Write([]byte(message))
		if err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}
		err = w.Close()
		if err != nil {
			return fmt.Errorf("failed to close writer: %w", err)
		}

		return client.Quit()
	}

	// For development (Ethereal, local SMTP without TLS)
	err := smtp.SendMail(addr, auth, s.From, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

// RenderTemplate renders an email template with data
func (s *EmailService) RenderTemplate(templateName string, data EmailData) (string, error) {
	// Set default app info
	data.AppName = s.AppConfig.Name
	data.AppURL = s.AppConfig.URL

	// Load template
	templatePath := filepath.Join("pkg", "email", "templates", templateName)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return body.String(), nil
}

// SendWelcomeEmail sends a welcome email to new users
func (s *EmailService) SendWelcomeEmail(name, email, verificationURL string) error {
	data := EmailData{
		Name:            name,
		Email:           email,
		VerificationURL: verificationURL,
	}

	htmlBody, err := s.RenderTemplate("welcome.html", data)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("Welcome to %s!", s.AppConfig.Name)
	return s.SendEmail(email, subject, htmlBody)
}

// SendVerificationEmail sends email verification link
func (s *EmailService) SendVerificationEmail(name, email, verificationURL string) error {
	data := EmailData{
		Name:            name,
		Email:           email,
		VerificationURL: verificationURL,
	}

	htmlBody, err := s.RenderTemplate("verifyEmail.html", data)
	if err != nil {
		return err
	}

	subject := "Verify Your Email Address"
	return s.SendEmail(email, subject, htmlBody)
}

// SendPasswordResetEmail sends password reset link
func (s *EmailService) SendPasswordResetEmail(name, email, resetURL string) error {
	data := EmailData{
		Name:     name,
		Email:    email,
		ResetURL: resetURL,
	}

	htmlBody, err := s.RenderTemplate("resetPassword.html", data)
	if err != nil {
		return err
	}

	subject := "Password Reset Request"
	return s.SendEmail(email, subject, htmlBody)
}

// SendPasswordChangedEmail sends confirmation after password change
func (s *EmailService) SendPasswordChangedEmail(name, email string) error {
	data := EmailData{
		Name:  name,
		Email: email,
	}

	htmlBody, err := s.RenderTemplate("passwordChanged.html", data)
	if err != nil {
		return err
	}

	subject := "Your Password Has Been Changed"
	return s.SendEmail(email, subject, htmlBody)
}
