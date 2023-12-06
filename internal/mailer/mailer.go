package mailer

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/thejixer/shop-api/internal/models"
)

type MailerService struct {
	from     string
	smtpHost string
	smtpPort string
	auth     smtp.Auth
}

func NewMailerService() *MailerService {

	from := os.Getenv("GMAIL_ADDRESS")
	password := os.Getenv("GMAIL_PASSWORD")
	smtpHost := "smtp.gmail.com"
	auth := smtp.PlainAuth("", from, password, smtpHost)

	return &MailerService{
		from:     os.Getenv("GMAIL_ADDRESS"),
		smtpHost: "smtp.gmail.com",
		smtpPort: "587",
		auth:     auth,
	}
}

func devPrintSkipEmail(env string) {
	fmt.Println("######################################################")
	fmt.Printf("skiping sending email since we're in %v enviroment \n", env)
	fmt.Println("######################################################")
}

func (m *MailerService) SendVerificationEmail(u *models.User, c string) error {
	env := os.Getenv("ENVIROMENT")
	if env == "DEV" || env == "TEST" {
		devPrintSkipEmail(env)
		return nil
	}

	DOMAIN := os.Getenv("DOMAIN")
	subject := "Subject: Verification Email from shop-api !\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
	<html>
		<body>
			<h2> Hello Dear %v, from shop-api </h2>
			<div>
				to activate your account, please click this 
				<a href="%v/auth/verify-email?email=%v&code=%v" > link </a>
				<br /><br />
				if you have <strong> not </strong> requested this, simply ignore this.
			</div>
		</body>
	</html>
	`, u.Name, DOMAIN, u.Email, c)
	msg := fmt.Sprintf("%v%v%v", subject, mime, body)
	message := []byte(msg)

	to := []string{
		u.Email,
	}

	addr := fmt.Sprintf("%v:%v", m.smtpHost, m.smtpPort)
	err := smtp.SendMail(addr, m.auth, m.from, to, message)
	if err != nil {
		return err
	}
	fmt.Println("Email Sent Successfully!")

	return nil
}

func (m *MailerService) SendPasswordChangeRequestEmail(u *models.User, c string) error {
	env := os.Getenv("ENVIROMENT")
	if env == "DEV" || env == "TEST" {
		devPrintSkipEmail(env)
		return nil
	}

	DOMAIN := os.Getenv("DOMAIN")
	subject := "Subject: change password request !\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<html>
			<body>
				<h2> Hello Dear %v </h2>
				<div>
					if you have requested to change your password, click this link
					<a href="%v/auth/verify-changepassword-request?email=%v&code=%v" > link </a> a<br /><br />
					if you have <strong> not </strong> requested this, simply ignore this.
				</div>
			</body>
		</html>
		`, u.Name, DOMAIN, u.Email, c)
	msg := fmt.Sprintf("%v%v%v", subject, mime, body)
	message := []byte(msg)

	to := []string{
		u.Email,
	}

	addr := fmt.Sprintf("%v:%v", m.smtpHost, m.smtpPort)
	err := smtp.SendMail(addr, m.auth, m.from, to, message)
	if err != nil {
		return err
	}
	fmt.Println("Email Sent Successfully!")
	return nil
}

func (m *MailerService) NotifyShipmentEmail(order *models.OrderDto) error {

	env := os.Getenv("ENVIROMENT")
	if env == "DEV" || env == "TEST" {
		devPrintSkipEmail(env)
		return nil
	}

	subject := "Subject: shipment notification !\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<html>
			<body>
				<h2> Hello Dear %v </h2>
				<div>
					your package has been sent and will be delivered soon to %v at the address below <br />
					<p> %v </p>
				</div>
			</body>
		</html>
		`, order.User.Name, order.Address.RecieverName, order.Address.Address)
	msg := fmt.Sprintf("%v%v%v", subject, mime, body)
	message := []byte(msg)

	to := []string{
		order.User.Email,
	}

	addr := fmt.Sprintf("%v:%v", m.smtpHost, m.smtpPort)
	err := smtp.SendMail(addr, m.auth, m.from, to, message)
	if err != nil {
		return err
	}
	fmt.Println("Email Sent Successfully!")
	return nil
}
