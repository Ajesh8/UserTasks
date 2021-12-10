package models

import (
	"errors"
	"fmt"
	"net/smtp"
	"time"

	"gorm.io/gorm"
)

type EmailAlert struct {
	gorm.Model
	Email     string `gorm:"" json:"email"`
	TaskId    uint   `json:"taskId"`
	AlertTime uint   `json:"alertTime"`
	Sent      bool   `json:"Sent"`
}

func (b *EmailAlert) CreateEmailAlert() (*EmailAlert, error) {
	if !validateEmailAlert(*b) {
		return nil, errors.New("email alert' mandatory fields not found")
	}
	b.Sent = false
	db.Create(&b)
	return b, nil
}

func validateEmailAlert(emailAlert EmailAlert) bool {
	if emailAlert.Email == "" {
		return false
	}
	if emailAlert.TaskId == 0 {
		return false
	}
	if emailAlert.AlertTime == 0 {
		return false
	}
	return true
}

func GetAllEmailAlert() []EmailAlert {
	var emailAlerts []EmailAlert
	db.Find(&emailAlerts)
	return emailAlerts
}

func AlertJob() {
	for {
		fmt.Println("Running email alert job")
		emailAlerts := GetAllEmailAlert()
		for _, emailAlert := range emailAlerts {
			task := FindTaskById(int64(emailAlert.TaskId))
			if time.Now().Add(time.Duration(emailAlert.AlertTime)*time.Hour).After(task.DueDate) && !emailAlert.Sent { //If current time + alert duration is greater than due date, send alert.
				err := SendEmailAlert(emailAlert, *task)
				if err != nil {
					fmt.Println("Error sending mail alert")
				} else {
					emailAlert.Sent = true
					db.Save(emailAlert)
				}
			}
		}
		time.Sleep(time.Minute * 5)
	}
}

func SendEmailAlert(emailAlert EmailAlert, task Task) error {
	from := ""
	password := ""

	to := []string{
		emailAlert.Email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	fromHeader := fmt.Sprintf("From: <%s>\r\n", from)
	toHeader := fmt.Sprintf("To: <%s>\r\n", to[0])
	subject := "Task Reminder"
	body := fmt.Sprintf("Reminder for your task with Title:%s and due date:%s", task.Title, task.DueDate.Format("2006-01-02 15:04:05"))
	msg := fromHeader + toHeader + subject + "\r\n" + body

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Email Sent Successfully!")
	return nil
}
