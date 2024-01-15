package effects

import "github.com/hajimehoshi/ebiten/v2"

const (
	WarningBubbleDuration = 60 * 5

	warningBubbleMaxScale = 2
)

type WarningBubble struct {
	tick    uint64
	avatar  *ebiten.Image
	message string

	x, y   float64
	cursor float64
}

func NewWarningBubble(x, y float64, message string, avatar *ebiten.Image) *WarningBubble {
	return &WarningBubble{
		tick:    WarningBubbleDuration,
		avatar:  avatar,
		message: message,

		x: x,
		y: y,
	}
}

func (wb *WarningBubble) Avatar() *ebiten.Image {
	return wb.avatar
}

func (wb *WarningBubble) Message() string {
	return wb.message
}

func (wb *WarningBubble) XY() (float64, float64) {
	return wb.x, wb.y
}

func (wb *WarningBubble) Cursor() float64 {
	return wb.cursor
}

func (wb *WarningBubble) Scale() float64 {
	return wb.cursor * warningBubbleMaxScale
}

func (wb *WarningBubble) Expired() bool {
	return wb == nil || wb.tick == 0
}

func (wb *WarningBubble) Update() {
	if wb.tick > 0 {
		wb.tick--
	}
	wb.cursor = 1 - float64(wb.tick)/WarningBubbleDuration
}
