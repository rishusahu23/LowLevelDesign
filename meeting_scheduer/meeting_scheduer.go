package main

import "time"

type RoomManager interface {
	BookRoom(roomId string, startTime, endTime time.Time, userId string) error
	CheckAvailability(roomId string, startTime, endTime time.Time) bool
	GetAvailableRooms(startTime, endTime time.Time) []Room
	GetMeetingHistory(userId string) []Meeting
}

type Notifier interface {
	NotifyUsers(meeting Meeting) error
	RegisterParticipant(participantId string)
	DeRegisterParticipant(participantId string)
}

type HistoryRepository interface {
	SaveMeeting(meeting Meeting) error
	GetMeetingByUser(userId string) []Meeting
	GetMeetingsForRoom(roomId string) []Meeting
}

type Meeting struct {
	Id           string
	RoomId       string
	StartTime    time.Time
	EndTime      time.Time
	Participants []string
	Organizer    string
}

type Room struct {
	Id       string
	Nane     string
	Capacity int
}

type User struct {
	Id    string
	Name  string
	Email string
}

type RoomManagerImpl struct {
	rooms      []Room
	notifier   Notifier
	repository HistoryRepository
}

var _ RoomManager = &RoomManagerImpl{}

func (r *RoomManagerImpl) BookRoom(roomId string, startTime, endTime time.Time, userId string) error {
	//TODO implement me
	panic("implement me")
}

func (r *RoomManagerImpl) CheckAvailability(roomId string, startTime, endTime time.Time) bool {
	meetings := r.repository.GetMeetingsForRoom(roomId)
	for _, meeting := range meetings {
		if isOverLapping(meeting.StartTime, meeting.EndTime, startTime, endTime) {
			return false
		}
	}
	return true
}

func (r *RoomManagerImpl) GetAvailableRooms(startTime, endTime time.Time) []Room {
	//TODO implement me
	panic("implement me")
}

func (r *RoomManagerImpl) GetMeetingHistory(userId string) []Meeting {
	return nil
}

func isOverLapping(existingStart, existingEnd, newStart, newEnd time.Time) bool {
	return !(newStart.After(existingEnd) || newEnd.Before(existingStart))
}
