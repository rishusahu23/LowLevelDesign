package main

import "fmt"

type Notifier interface {
	Send(message string)
}

type EmailNotifier struct {
}

func (e *EmailNotifier) Send(message string) {
	fmt.Println("sending email with message: ", message)
}

type SMSDecorator struct {
	notifier Notifier
}

func (s *SMSDecorator) Send(message string) {
	s.notifier.Send(message)
	fmt.Println("sending sms with notifier: ", message)
}

func main() {
	notifier := &EmailNotifier{}
	notifier.Send("email notifier")

	smsDecorator := &SMSDecorator{
		notifier: notifier,
	}
	smsDecorator.Send("sms decorator")
}