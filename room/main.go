package main

import "fmt"

type MailClient struct {
}

func (m *MailClient) SendMail(email string) {
	fmt.Println("sending mail")
}

type Interval struct {
	StartTime int64
	EndTime   int64
}

type User struct {
	Email, name string
}

type Meeting struct {
	Id        string
	StartTime int64
	EndTime   int64
	Invitees  []User
}

type MeetingRoom struct {
	Id         string
	Capacity   int
	Calender   Calender
	MailClient MailClient
}

type Calender struct {
	MeetingList []Meeting
	Calender    map[int64]int64
}

func (c *Calender) CheckAvailability(interval Interval) bool {
	for _, meeting := range c.MeetingList {
		if (interval.StartTime >= meeting.StartTime && interval.StartTime < meeting.EndTime) ||
			(interval.EndTime > meeting.StartTime && interval.EndTime <= meeting.EndTime) ||
			(meeting.StartTime >= interval.StartTime && meeting.EndTime <= interval.EndTime) {
			return false
		}
	}
	return true
}
