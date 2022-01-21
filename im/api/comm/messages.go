package comm

type NewContactMessage struct {
	Id       int64
	Type     int
	FromId   int64
	FromType int
}
