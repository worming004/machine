package bubblemachine

type Bubble string

func (b Bubble) String() string {
	return string(b)
}

func NewBubble(name string) Bubble {
	return Bubble(name)
}

func printableBubbles(b []Bubble) string {
	var s string
	for _, v := range b {
		s += string(v) + ", "
	}
	return "[" + s[:len(s)-2] + "]"
}
