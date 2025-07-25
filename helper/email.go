package helper

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"go-face-auth/config"
	"log"
	"net/smtp"
	"net/textproto"
	
)

// SendPaymentLinkEmail sends an email with the payment link to the admin.
func SendPaymentLinkEmail(recipientEmail, companyName, paymentLink string) error {
	// Check if SMTP configuration is loaded
	if config.SMTP_SERVER == "" || config.SMTP_PORT == "" || config.SMTP_USER == "" || config.SMTP_PASSWORD == "" || config.SMTP_FROM == "" {
		log.Println("Skipping email sending: SMTP configuration is incomplete.")
		return fmt.Errorf("SMTP configuration incomplete")
	}

	subject := fmt.Sprintf("Pembayaran Langganan %s Anda", companyName)
	body := fmt.Sprintf(`
		<h1>Pembayaran Langganan Anda</h1>
		<p>Halo Admin %s,</p>

		<p>Terima kasih telah mendaftar untuk langganan %s. Silakan klik tautan di bawah ini untuk menyelesaikan pembayaran Anda:</p>
		<p><a href="%s">Lanjutkan Pembayaran</a></p>
		<p>Jika Anda memiliki pertanyaan, jangan ragu untuk menghubungi kami.</p>
		<p>Hormat kami,<br>Tim Go-Face-Auth</p>
	`, companyName, companyName, paymentLink)

	message := []byte(
		"To: " + recipientEmail + "\r\n" +
		"From: " + config.SMTP_FROM + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body,
	)

	return sendMail(recipientEmail, message)
}

// SendInvoiceEmail sends an email with an attached PDF invoice.
func SendInvoiceEmail(recipientEmail, companyName, invoiceFileName string, invoicePDFData []byte) error {
	// Check if SMTP configuration is loaded
	if config.SMTP_SERVER == "" || config.SMTP_PORT == "" || config.SMTP_USER == "" || config.SMTP_PASSWORD == "" || config.SMTP_FROM == "" {
		log.Println("Skipping email sending: SMTP configuration is incomplete.")
		return fmt.Errorf("SMTP configuration incomplete")
	}

	subject := fmt.Sprintf("Invoice Pembayaran Langganan %s", companyName)

	// Email body (HTML part)
	htmlBody := fmt.Sprintf(`
		<h1>Invoice Pembayaran Anda</h1>
		<p>Halo Admin %s,</p>
		<p>Terima kasih atas pembayaran langganan Anda. Terlampir adalah invoice pembayaran Anda.</p>
		<p>Jika Anda memiliki pertanyaan, jangan ragu untuk menghubungi kami.</p>
		<p>Hormat kami,<br>Tim Go-Face-Auth</p>
	`, companyName)

	// Create a new buffer for the email message
	buf := new(bytes.Buffer)
	bw := bufio.NewWriter(buf)
	w := textproto.NewWriter(bw)

	// Write email headers
	fmt.Fprintf(w.W, "From: %s\r\n", config.SMTP_FROM)
	fmt.Fprintf(w.W, "To: %s\r\n", recipientEmail)
	fmt.Fprintf(w.W, "Subject: %s\r\n", subject)
	fmt.Fprintf(w.W, "MIME-Version: 1.0\r\n")

	// Generate a random boundary string
	boundary := "GoBoundary" // A simple boundary, for more robust, use a UUID
	fmt.Fprintf(w.W, "Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary)
	fmt.Fprintf(w.W, "\r\n") // End of headers

	// Write HTML part
	fmt.Fprintf(w.W, "--%s\r\n", boundary)
	fmt.Fprintf(w.W, "Content-Type: text/html; charset=\"UTF-8\"\r\n")
	fmt.Fprintf(w.W, "Content-Transfer-Encoding: quoted-printable\r\n")
	fmt.Fprintf(w.W, "\r\n")
	fmt.Fprintf(w.W, "%s", htmlBody)
	fmt.Fprintf(w.W, "\r\n")

	// Write PDF attachment part
	fmt.Fprintf(w.W, "--%s\r\n", boundary)
	fmt.Fprintf(w.W, "Content-Type: application/pdf\r\n")
	fmt.Fprintf(w.W, "Content-Transfer-Encoding: base64\r\n")
	fmt.Fprintf(w.W, "Content-Disposition: attachment; filename=\"%s\"\r\n", invoiceFileName)
	fmt.Fprintf(w.W, "\r\n")

	encoder := base64.NewEncoder(base64.StdEncoding, w.W)
	encoder.Write(invoicePDFData)
	encoder.Close()
	fmt.Fprintf(w.W, "\r\n")

	// End of multipart message
	fmt.Fprintf(w.W, "--%s--\r\n", boundary)
	bw.Flush()

	return sendMail(recipientEmail, buf.Bytes())
}

// SendPasswordResetEmail sends an email with a password reset link.
func SendPasswordResetEmail(recipientEmail, recipientName, resetURL string) error {
	// Check if SMTP configuration is loaded
	if config.SMTP_SERVER == "" || config.SMTP_PORT == "" || config.SMTP_USER == "" || config.SMTP_PASSWORD == "" || config.SMTP_FROM == "" {
		log.Println("Skipping email sending: SMTP configuration is incomplete.")
		return fmt.Errorf("SMTP configuration incomplete")
	}

	subject := "Atur Kata Sandi Akun Karyawan Anda"
	body := fmt.Sprintf(`
		<h1>Atur Kata Sandi Anda</h1>
		<p>Halo %s,</p>
		<p>Akun karyawan Anda telah dibuat. Silakan klik tautan di bawah ini untuk mengatur kata sandi Anda:</p>
		<p><a href="%s">Atur Kata Sandi</a></p>
		<p>Tautan ini akan kedaluwarsa dalam 24 jam.</p>
		<p>Jika Anda tidak meminta ini, abaikan email ini.</p>
		<p>Hormat kami,<br>Tim Go-Face-Auth</p>
	`, recipientName, resetURL)

	message := []byte(
		"To: " + recipientEmail + "\r\n" +
		"From: " + config.SMTP_FROM + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body,
	)

	return sendMail(recipientEmail, message)
}

// SendSubscriptionReminderEmail sends a reminder email for an upcoming subscription renewal.
func SendSubscriptionReminderEmail(recipientEmail, companyName string, daysRemaining int, renewalLink string) error {
	// Check if SMTP configuration is loaded
	if config.SMTP_SERVER == "" || config.SMTP_PORT == "" || config.SMTP_USER == "" || config.SMTP_PASSWORD == "" || config.SMTP_FROM == "" {
		log.Println("Skipping email sending: SMTP configuration is incomplete.")
		return fmt.Errorf("SMTP configuration incomplete")
	}

	subject := fmt.Sprintf("Pemberitahuan Perpanjangan Langganan %s Anda", companyName)
	body := fmt.Sprintf(`
		<h1>Pemberitahuan Perpanjangan Langganan</h1>
		<p>Halo Admin %s,</p>
		<p>Langganan %s Anda akan segera berakhir dalam <b>%d hari</b>.</p>
		<p>Untuk memastikan layanan tidak terganggu, silakan perpanjang langganan Anda sekarang:</p>
		<p><a href="%s">Perpanjang Langganan Sekarang</a></p>
		<p>Terima kasih atas kepercayaan Anda.</p>
		<p>Hormat kami,<br>Tim Go-Face-Auth</p>
	`, companyName, companyName, daysRemaining, renewalLink)

	message := []byte(
		"To: " + recipientEmail + "\r\n" +
		"From: " + config.SMTP_FROM + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body,
	)

	return sendMail(recipientEmail, message)
}

// SendSubscriptionExpiredEmail sends an email notification that a subscription has expired.
func SendSubscriptionExpiredEmail(recipientEmail, companyName string, renewalLink string) error {
	// Check if SMTP configuration is loaded
	if config.SMTP_SERVER == "" || config.SMTP_PORT == "" || config.SMTP_USER == "" || config.SMTP_PASSWORD == "" || config.SMTP_FROM == "" {
		log.Println("Skipping email sending: SMTP configuration is incomplete.")
		return fmt.Errorf("SMTP configuration incomplete")
	}

	subject := fmt.Sprintf("Langganan %s Anda Telah Kedaluwarsa", companyName)
	body := fmt.Sprintf(`
		<h1>Langganan Anda Telah Kedaluwarsa</h1>
		<p>Halo Admin %s,</p>
		<p>Langganan %s Anda telah kedaluwarsa. Beberapa fitur mungkin tidak lagi dapat diakses.</p>
		<p>Untuk mengaktifkan kembali layanan, silakan perpanjang langganan Anda:</p>
		<p><a href="%s">Perpanjangan Langganan Sekarang</a></p>
		<p>Terima kasih atas pengertian Anda.</p>
		<p>Hormat kami,<br>Tim Go-Face-Auth</p>
	`, companyName, companyName, renewalLink)

	message := []byte(
		"To: " + recipientEmail + "\r\n" +
		"From: " + config.SMTP_FROM + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body,
	)

	return sendMail(recipientEmail, message)
}

// SendConfirmationEmail sends an email with a confirmation link for new company admin registrations.
func SendConfirmationEmail(recipientEmail, companyName, confirmationLink string) error {
	// Check if SMTP configuration is loaded
	if config.SMTP_SERVER == "" || config.SMTP_PORT == "" || config.SMTP_USER == "" || config.SMTP_PASSWORD == "" || config.SMTP_FROM == "" {
		log.Println("Skipping email sending: SMTP configuration is incomplete.")
		return fmt.Errorf("SMTP configuration incomplete")
	}

	subject := fmt.Sprintf("Konfirmasi Pendaftaran Perusahaan %s Anda", companyName)
	body := fmt.Sprintf(`
		<h1>Konfirmasi Pendaftaran Anda</h1>
		<p>Halo Admin %s,</p>
		<p>Terima kasih telah mendaftar untuk %s. Silakan klik tautan di bawah ini untuk mengkonfirmasi alamat email Anda dan mengaktifkan akun Anda:</p>
		<p><a href="%s">Konfirmasi Email Anda</a></p>
		<p>Jika Anda tidak mendaftar untuk ini, abaikan email ini.</p>
		<p>Hormat kami,<br>Tim Go-Face-Auth</p>
	`, companyName, companyName, confirmationLink)

	message := []byte(
		"To: " + recipientEmail + "\r\n" +
		"From: " + config.SMTP_FROM + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body,
	)

	return sendMail(recipientEmail, message)
}

// GetFrontendBaseURL returns the configured frontend base URL.
func GetFrontendBaseURL() string {
	return config.FrontendBaseURL
}

// GetFrontendAdminBaseURL returns the configured admin frontend base URL.
func GetFrontendAdminBaseURL() string {
	return config.FrontendAdminBaseURL
}

// sendMail is a helper function to handle the actual SMTP sending logic.
func sendMail(recipientEmail string, message []byte) error {
	auth := smtp.PlainAuth("", config.SMTP_USER, config.SMTP_PASSWORD, config.SMTP_SERVER)

	conn, err := smtp.Dial(config.SMTP_SERVER + ":" + config.SMTP_PORT)
	if err != nil {
		log.Printf("Error connecting to SMTP server: %v", err)
		return err
	}
	defer conn.Close()

	if err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
		log.Printf("Error starting TLS: %v", err)
		return err
	}

	if err = conn.Auth(auth); err != nil {
		log.Printf("Error authenticating with SMTP server: %v", err)
		return err
	}

	if err = conn.Mail(config.SMTP_FROM); err != nil {
		log.Printf("Error setting sender email: %v", err)
		return err
	}
	if err = conn.Rcpt(recipientEmail); err != nil {
		log.Printf("Error setting recipient email: %v", err)
		return err
	}

	w, err := conn.Data()
	if err != nil {
		log.Printf("Error getting data writer: %v", err)
		return err
	}
	_, err = w.Write(message)
	if err != nil {
		log.Printf("Error writing email message: %v", err)
		return err
	}
	err = w.Close()
	if err != nil {
		log.Printf("Error closing data writer: %v", err)
		return err
	}

	log.Printf("Email sent to %s.", recipientEmail)
	return nil
}

