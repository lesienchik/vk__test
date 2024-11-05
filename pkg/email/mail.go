package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/lesienchik/vk__test/internal/config"
)

type Email struct {
	addr     string
	password string
	site     string
}

func New(cfg *config.Email) *Email {
	return &Email{
		addr:     cfg.Addr,
		password: cfg.Password,
		site:     cfg.Site,
	}
}

// Отправляет код подтверждения на почту пользователя, чтобы тот мог завершить регистрацию.
func (m *Email) SendConfirmCode(to, verifyCode string) error {
	// Настройка SMTP клиента.
	var (
		smtpHost = "smtp.mail.ru"
		smtpPort = "465"
	)

	// Настройка TLS.
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	})
	if err != nil {
		return fmt.Errorf("email.SendConfirmCode(1): %w", err)
	}

	smtpClient, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("email.SendConfirmCode(2): %w", err)
	}

	// Аутентификация.
	auth := smtp.PlainAuth("", m.addr, m.password, smtpHost)
	if err := smtpClient.Auth(auth); err != nil {
		return fmt.Errorf("email.SendConfirmCode(3): %w", err)
	}

	// Установка отправителя и получателя.
	if err := smtpClient.Mail(m.addr); err != nil {
		return fmt.Errorf("email.SendConfirmCode(4): %w", err)
	}
	if err := smtpClient.Rcpt(to); err != nil {
		return fmt.Errorf("email.SendConfirmCode(5): %w", err)
	}

	// Формирование сообщения.
	message := m.getConfirmMessage(to, verifyCode)

	// Отправка сообщения
	w, err := smtpClient.Data()
	if err != nil {
		return fmt.Errorf("email.SendConfirmCode(6): %w", err)
	}

	if _, err := w.Write(message); err != nil {
		return fmt.Errorf("email.SendConfirmCode(7): %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("email.SendConfirmCode(8): %w", err)
	}
	return smtpClient.Quit()
}

// Формирует сообщение для завершения регистрации пользователя.
func (m *Email) getConfirmMessage(to, verifyCode string) []byte {
	link := fmt.Sprintf("%s/verify?code=%s", m.site, verifyCode)
	subject := "Subject: Подтверждение регистрации в vktest\n"
	body := fmt.Sprintf(`Здравствуйте, %s!

Рады приветствовать Вас в vktest.

Чтобы завершить регистрацию, пожалуйста, подтвердите адрес электронной почты, перейдя по ссылке ниже:

%s

С наилучшими пожеланиями,
Команда vktest`, to, link)
	return []byte(subject + "\n" + body)
}
