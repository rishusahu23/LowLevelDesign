package main

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	name    string
	id      string
	isAdmin bool
}

type Plan int32

const (
	ONE_MONTH Plan = iota
	THREE_MONTH
	SIX_MONTHS
	TWELVE_MONTHS
)

type Subscription struct {
	id        string
	plan      Plan
	expiresAt time.Time
	userId    string
}

type GameType int32

const (
	CRICKET GameType = iota
	FOOTBALL
)

type IScore interface {
	UpdateScore(newScore string)
	GetScore() string
}

type FootballScore struct {
	score string
}

func (f *FootballScore) UpdateScore(newScore string) {
	f.score = newScore
}

func (f *FootballScore) GetScore() string {
	return f.score
}

func NewFootballScore() *FootballScore {
	return &FootballScore{}
}

type CricketScore struct {
	score string
}

func (f *CricketScore) UpdateScore(newScore string) {
	f.score = newScore
}

func (f *CricketScore) GetScore() string {
	return f.score
}

func NewCricketScore() *CricketScore {
	return &CricketScore{}
}

type Game struct {
	id       string
	gameType GameType
	score    IScore
	history  string
	status   string
}

type IUserSvc interface {
	Register(user *User) error
	GetById(id string) (*User, error)
}

type ISubscription interface {
	Add(subscription *Subscription) error
	GetById(id string) (*Subscription, error)
}

type IGame interface {
	Add(game *Game) error
	GetById(id string) (*Game, error)
	UpdateScore(id, score string) error
	GetHistory(id string) (string, error)
}

type UserSvc struct {
	users map[string]*User
}

func NewUserSvc() *UserSvc {
	return &UserSvc{
		users: make(map[string]*User),
	}
}

func (u *UserSvc) Register(user *User) error {
	u.users[user.id] = user
	return nil
}

func (u *UserSvc) GetById(id string) (*User, error) {
	user, ok := u.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return user, nil
}

type SubscriptionSvc struct {
	subs map[string]*Subscription
}

func NewSubscriptionSvc() *SubscriptionSvc {
	return &SubscriptionSvc{
		subs: make(map[string]*Subscription),
	}
}

func (s *SubscriptionSvc) Add(subscription *Subscription) error {
	s.subs[subscription.id] = subscription
	return nil
}

func (s *SubscriptionSvc) GetById(id string) (*Subscription, error) {
	sub, ok := s.subs[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return sub, nil
}

type GameSvc struct {
	games map[string]*Game
}

func NewGameSvc() *GameSvc {
	return &GameSvc{games: make(map[string]*Game)}
}

func (g *GameSvc) Add(game *Game) error {
	g.games[game.id] = game
	return nil

}

func (g *GameSvc) GetById(id string) (*Game, error) {
	game, exists := g.games[id]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}
	return game, nil
}

func (g *GameSvc) UpdateScore(id, score string) error {
	game, exists := g.games[id]
	if !exists {
		return fmt.Errorf("game not found")
	}
	game.score.UpdateScore(score)
	g.games[id] = game
	return nil
}

func (g *GameSvc) GetHistory(id string) (string, error) {
	game, exists := g.games[id]
	if !exists {
		return "", fmt.Errorf("game not found")
	}
	return game.history, nil
}

type Service struct {
	userSvc IUserSvc
	gameSvc IGame
	subSvc  ISubscription
}

func NewService() *Service {
	return &Service{
		userSvc: NewUserSvc(),
		gameSvc: NewGameSvc(),
		subSvc:  NewSubscriptionSvc(),
	}
}

func (s *Service) RegisterUser(user *User) error {
	return s.userSvc.Register(user)
}

func (s *Service) GetUserByID(userID string) (*User, error) {
	return s.userSvc.GetById(userID)
}

func (s *Service) AddSubscription(sub *Subscription) error {
	return s.subSvc.Add(sub)
}

func (s *Service) GetSubscriptionByUserID(userID string) (*Subscription, error) {
	return s.subSvc.GetById(userID)
}

func (s *Service) AddGame(game *Game) error {
	return s.gameSvc.Add(game)
}

func (s *Service) GetGameByID(gameID string) (*Game, error) {
	return s.gameSvc.GetById(gameID)
}

func (s *Service) UpdateGameScore(gameID string, score string) error {
	return s.gameSvc.UpdateScore(gameID, score)
}

func (s *Service) GetGameHistory(gameID string) (string, error) {
	return s.gameSvc.GetHistory(gameID)
}

type GameFactory struct {
}

func (gf *GameFactory) CreateGame(id string, gameType GameType) *Game {
	var score IScore
	switch gameType {
	case CRICKET:
		score = &CricketScore{}
	case FOOTBALL:
		score = &FootballScore{}
	}
	return &Game{
		id:       id,
		gameType: gameType,
		score:    score,
		history:  "none",
		status:   "Pending",
	}
}

func main() {

}
