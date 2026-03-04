package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendPasswordResetEmail sends a password reset email
func SendPasswordResetEmail(toEmail, resetLink string) error {
	// Get email configuration from environment
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SENDER_EMAIL")
	senderPassword := os.Getenv("SENDER_PASSWORD")

	// Default values for development
	if smtpHost == "" {
		smtpHost = "localhost"
	}
	if smtpPort == "" {
		smtpPort = "1025" // Mailhog default
	}
	if senderEmail == "" {
		senderEmail = "noreply@unieroll.local"
	}

	// Email content
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.button { background-color: #007bff; color: white; padding: 10px 20px; border-radius: 5px; text-decoration: none; display: inline-block; margin: 20px 0; }
	</style>
</head>
<body>
	<div class="container">
		<h2>Password Reset Request</h2>
		<p>You requested a password reset for your UniEnroll account.</p>
		<p>Click the button below to reset your password. This link will expire in 15 minutes.</p>
		<a href="%s" class="button">Reset Password</a>
		<p>If you did not request this, please ignore this email.</p>
		<hr>
		<p><small>UniEnroll Team</small></p>
	</div>
</body>
</html>
`, resetLink)

	// Compose email
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n", senderEmail, toEmail, subject)
	message := headers + body

	// For development, if password is not set, skip SMTP sending
	if senderPassword == "" {
		fmt.Printf("📧 [DEV] Password reset email would be sent to: %s\n", toEmail)
		fmt.Printf("   Reset Link: %s\n", resetLink)
		return nil
	}

	// Send email via SMTP
	addr := smtpHost + ":" + smtpPort
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	err := smtp.SendMail(addr, auth, senderEmail, []string{toEmail}, []byte(message))
	if err != nil {
		// Log error but don't fail in development
		fmt.Printf("⚠️ Failed to send email: %v\n", err)
		return nil
	}

	return nil
}

// SendWelcomeEmail sends a welcome email with login credentials
func SendWelcomeEmail(toEmail, name, userID, password, role string) error {
	// Get email configuration from environment
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SENDER_EMAIL")
	senderPassword := os.Getenv("SENDER_PASSWORD")

	// Default values for development
	if smtpHost == "" {
		smtpHost = "localhost"
	}
	if smtpPort == "" {
		smtpPort = "1025" // Mailhog default
	}
	if senderEmail == "" {
		senderEmail = "noreply@unieroll.local"
	}

	// Email content
	subject := "Welcome to UniEnroll - Your Account Details"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f9f9f9; border-radius: 10px; }
		.header { background-color: #007bff; color: white; padding: 20px; border-radius: 10px 10px 0 0; text-align: center; }
		.content { background-color: white; padding: 30px; border-radius: 0 0 10px 10px; }
		.credentials { background-color: #f0f8ff; padding: 15px; border-left: 4px solid #007bff; margin: 20px 0; }
		.button { background-color: #28a745; color: white; padding: 12px 30px; border-radius: 5px; text-decoration: none; display: inline-block; margin: 20px 0; }
		.warning { color: #d9534f; font-weight: bold; margin: 15px 0; }
		.footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>🎓 Welcome to UniEnroll</h1>
		</div>
		<div class="content">
			<h2>Hello %s,</h2>
			<p>Your account has been created successfully by the administrator.</p>
			
			<div class="credentials">
				<h3>📋 Your Login Credentials:</h3>
				<p><strong>Email:</strong> %s</p>
				<p><strong>Temporary Password:</strong> <code>%s</code></p>
				<p><strong>Role:</strong> %s</p>
			</div>
			
			<p class="warning">⚠️ Important: Please change your password after your first login for security purposes.</p>
			
			<p>You can now access the UniEnroll system using the credentials above.</p>
			
			<a href="http://localhost:3000/login" class="button">Login to UniEnroll</a>
			
			<hr style="margin: 30px 0; border: none; border-top: 1px solid #ddd;">
			
			<h3>📚 What's Next?</h3>
			<ul>
				<li><strong>Students:</strong> View and enroll in available modules</li>
				<li><strong>Lecturers:</strong> Access your assigned modules</li>
				<li><strong>Admins:</strong> Manage the entire system</li>
			</ul>
			
			<p>If you have any questions or need assistance, please contact the administrator.</p>
			
			<div class="footer">
				<p>UniEnroll System</p>
				<p><small>This is an automated message. Please do not reply to this email.</small></p>
			</div>
		</div>
	</div>
</body>
</html>
`, name, toEmail, password, role)

	// Compose email
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n", senderEmail, toEmail, subject)
	message := headers + body

	// For development, if password is not set, skip SMTP sending but show the credentials
	if senderPassword == "" {
		fmt.Printf("\n📧 [DEV] Welcome email would be sent to: %s\n", toEmail)
		fmt.Printf("   Name: %s\n", name)
		fmt.Printf("   Email: %s\n", toEmail)
		fmt.Printf("   Password: %s\n", password)
		fmt.Printf("   Role: %s\n\n", role)
		return nil
	}

	// Send email via SMTP
	addr := smtpHost + ":" + smtpPort
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	err := smtp.SendMail(addr, auth, senderEmail, []string{toEmail}, []byte(message))
	if err != nil {
		// Log error but don't fail in development
		fmt.Printf("⚠️ Failed to send email: %v\n", err)
		fmt.Printf("   Credentials: %s / %s\n", toEmail, password)
		return nil
	}

	fmt.Printf("✉️ Welcome email sent to: %s\n", toEmail)
	return nil
}
