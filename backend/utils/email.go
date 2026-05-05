package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bianquiviri/control-horario/models"
	"gopkg.in/gomail.v2"
)

func SendNotificationEmail(to string, log models.TimeLog) {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")

	if host == "" || user == "" {
		fmt.Println("SMTP not configured, skipping email")
		return
	}

	port, _ := strconv.Atoi(portStr)

	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Registro de Horas - Control de Horario")
	
	body := fmt.Sprintf(`
		<h1>Registro de Horas Finalizado</h1>
		<p>Has registrado una nueva sesión:</p>
		<ul>
			<li><b>Inicio:</b> %s</li>
			<li><b>Fin:</b> %s</li>
			<li><b>Duración:</b> %.2f horas</li>
		</ul>
	`, log.StartTime.Format("02/01/2006 15:04"), log.EndTime.Format("02/01/2006 15:04"), log.Duration)
	
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, user, pass)

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Could not send email: %v\n", err)
	}
}
