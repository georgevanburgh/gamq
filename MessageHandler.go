package gamq

type MessageHandler interface {
	Initialize(<-chan string) chan<- string
}
