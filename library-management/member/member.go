package member

type Member struct {
	Id           string
	Name         string
	MaxBooks     int
	CurrentBooks map[string]string
	Barcode      string
}
