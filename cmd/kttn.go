package main

import (
	"fmt"
	"math/rand/v2"
  "math"
	"strings"
  "unicode/utf8"

	"github.com/gdamore/tcell/v2"
)

const LOGO string = `
   /\_/\
   >^¸^<
    /|\
   (_|_)   `

const HAPPY_EYE string = "^"
const WATCHING_EYE string = "•"
const SURPRISED_EYE string = "°"
const SLEEPY_EYE string = "¯"

const TAIL_BEHIND string = "   (_|_)"
const TAIL_LEFT_FLAT string = "___(_|_)"
const TAIL_RIGHT_FLAT string = "   (_|_)___"
const TAIL_LEFT_UP string = "\\__(_|_)"
const TAIL_RIGHT_UP string = "   (_|_)__/"

const DEFAULT_TEXT string = 
  `This is just the default text for example and testing purposes. There is not point to this other than that. This is all just a very long string with no line breaks to illustrate some challenges with presenting it in a box and reflowing it when the window is resized.`

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Printf("%+v", err)
	}
	if err := s.Init(); err != nil {
		fmt.Printf("%+v", err)
	}

	// Set default text style
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	// Clear screen
	s.Clear()
  
	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
    s.Clear()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

  cursorPos := 0

	for {
    // Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
      s.Clear()
      updateLogo(s, defStyle)
      updateTypingBox(s, defStyle, DEFAULT_TEXT)
      updateCursor(s, DEFAULT_TEXT, cursorPos)
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
        cursorPos--
        updateCursor(s, DEFAULT_TEXT, cursorPos)
			} else {
        //TODO once text len == cursor pos stop the increment and show results
        cursorPos++
        updateLogo(s, defStyle)
        updateTypingBox(s, defStyle, DEFAULT_TEXT)
        updateCursor(s, DEFAULT_TEXT, cursorPos)
      }
		}
	}
}


func getMidScreenCoords(screen tcell.Screen) (midX int, midY int) {
  availX, availY := screen.Size()
  return availX/2, availY/2
}

func getLogoCoords(screen tcell.Screen) (logoX, logoY int) {
  midX, midY := getMidScreenCoords(screen)
  return midX - 6, midY - 8 
}

func getTypingBoxCoords(screen tcell.Screen, text string) (startX, startY, endX, endY int) {
  maxX, _ := screen.Size()
  _, midY := getMidScreenCoords(screen) 
  txtLength := utf8.RuneCountInString(text)
  //add .5 to always round up
  numOfRows := math.Round(float64(txtLength)/float64(maxX-12) + 0.5)  
  offset := 0
  if numOfRows > 6 {
    offset = int(numOfRows) - 6
  }
  return 5, midY - int(math.Round(numOfRows/2 + 0.5)) + offset, maxX - 5, midY + int(math.Round(numOfRows/2 + 0.5)) + offset
}

func updateCursor(screen tcell.Screen, text string, cursorPos int) {
  sx, sy, ex, _ := getTypingBoxCoords(screen, text)
  screen.SetCursorStyle(tcell.CursorStyleBlinkingBlock)
  lineLen := ex - sx - 1 
  yOffset := cursorPos/lineLen + 1
  xOffset := 0 
  if cursorPos >= lineLen {
    xOffset = lineLen * (cursorPos/lineLen) 
  }
  screen.ShowCursor(sx + 1 + cursorPos - xOffset, sy + yOffset)
}

func updateLogo(screen tcell.Screen, style tcell.Style) {
  x, y := getLogoCoords(screen)
  n := rand.IntN(5)
  switch n {
    case 0:
      drawText(screen, x, y, style, LOGO)
    case 1:
      drawText(screen, x, y, style, getLogoWithParams(WATCHING_EYE, TAIL_RIGHT_FLAT))
    case 2:
      drawText(screen, x, y, style, getLogoWithParams(HAPPY_EYE, TAIL_LEFT_FLAT))
    case 3:
      drawText(screen, x, y, style, getLogoWithParams(SURPRISED_EYE, TAIL_RIGHT_UP))
    case 4:
      drawText(screen, x, y, style, getLogoWithParams(SLEEPY_EYE, TAIL_LEFT_UP))
    
  }
}

func updateTypingBox(screen tcell.Screen, style tcell.Style, text string) {
  sx, sy, ex,ey := getTypingBoxCoords(screen, text) 
  drawBox(screen, sx, sy, ex, ey, style, text)
}

func getLogoWithParams(eyes, tail string) string {
  r := strings.Replace(LOGO, "^", eyes, 2)
  return strings.Replace(r, TAIL_BEHIND, tail, 1)
}

func drawBoundedText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func drawText(s tcell.Screen, x, y int, style tcell.Style, text string) {
	row := y
	col := x
  lines := strings.Split(text, "\n")
  for _, line := range lines {
    for _, r := range []rune(line) {
      s.SetContent(col, row, r, nil, style)
      col++
    }
    col = x
    row++
  }
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, '╭', nil, style)
		s.SetContent(x1, y2, '╰', nil, style)
		s.SetContent(x2, y1, '╮', nil, style)
		s.SetContent(x2, y2, '╯', nil, style)
	}

	drawBoundedText(s, x1+1, y1+1, x2, y2-1, style, text)
}

