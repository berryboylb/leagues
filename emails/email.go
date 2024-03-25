package emails

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"

	// "github.com/joho/godotenv"
)

var Auth smtp.Auth
var user string

func init() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	user = os.Getenv("SMTP_USER") // Removed := operator
	if user == "" {
		log.Fatal("Error loading smtp user variables")
	}

	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
		log.Fatal("Error loading smtp password variables")
	}
	Auth = smtp.PlainAuth("", user, smtpPassword, "smtp.gmail.com")
}

type EmailData struct {
	OTP string
}

func SendOTPEmail(userEmail string, otp string, title string) error {
	tmpl, err := template.ParseFiles("emails/templates/otp.html")
	if err != nil {
		return err
	}

	data := EmailData{
		OTP: otp,
	}

	var tpl bytes.Buffer
	if err = tmpl.Execute(&tpl, data); err != nil {
		return err
	}

	headers := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"From: GoStoreApp <phemmynesce4life@gmail.com>\r\n" +
		"To: " + userEmail + "\r\n" +
		"Subject: " + title + " \r\n"

	msg := []byte(headers + "\r\n" + tpl.String())

	fmt.Println("start sending email")
	err = smtp.SendMail("smtp.gmail.com:587", Auth, user, []string{userEmail}, msg)
	fmt.Println("stop sending email")
	return err
}
