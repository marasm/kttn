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

type Cat struct {
  CurEyes string
  CurTail string
}

type TypingTest struct {
  Text string
  Results []bool 
  CurPos int
}

func (t *TypingTest) UpdateWithRegKey(key rune) {
  // - see if the key == the current rune in text
  // - append the results with true|false
  t.Results[t.CurPos] = []rune(t.Text)[t.CurPos] == key
  if t.CurPos < utf8.RuneCountInString(t.Text) - 1 {
    t.CurPos++
  }
}

func (t *TypingTest) UpdateWithBackspace(key tcell.Key) {
  if t.CurPos > 0 {
    t.CurPos--
  }
}
  
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
  s.SetCursorStyle(tcell.CursorStyleBlinkingBlock)

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

  typeTest := TypingTest{
    Text: DEFAULT_TEXT,
    CurPos: 0,
    Results: make([]bool, utf8.RuneCountInString(DEFAULT_TEXT)),
  } 

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
      updateTypingBox(s, defStyle, typeTest)
      updateCursor(s, typeTest)
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
        typeTest.UpdateWithBackspace(ev.Key())
        updateTypingBox(s, defStyle, typeTest)
        updateCursor(s, typeTest)
			} else {
        //TODO once text len == cursor show results
        typeTest.UpdateWithRegKey(ev.Rune())
        updateLogo(s, defStyle)
        updateTypingBox(s, defStyle, typeTest)
        updateCursor(s, typeTest)
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

func updateCursor(screen tcell.Screen, typeTest TypingTest) {
  sx, sy, ex, _ := getTypingBoxCoords(screen, typeTest.Text)
  lineLen := ex - sx - 1 
  yOffset := typeTest.CurPos/lineLen + 1
  xOffset := 0 
  if typeTest.CurPos >= lineLen {
    xOffset = lineLen * (typeTest.CurPos/lineLen) 
  }
  screen.ShowCursor(sx + 1 + typeTest.CurPos - xOffset, sy + yOffset)
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

func updateTypingBox(screen tcell.Screen, style tcell.Style, typeTest TypingTest) {
  sx, sy, ex,ey := getTypingBoxCoords(screen, typeTest.Text) 
  drawBox(screen, sx, sy, ex, ey, style)
	drawBoundedText(screen, sx+1, sy+1, ex, ey-1, style, typeTest)
}

func getLogoWithParams(eyes, tail string) string {
  r := strings.Replace(LOGO, "^", eyes, 2)
  return strings.Replace(r, TAIL_BEHIND, tail, 1)
}

func drawBoundedText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, typeTest TypingTest) {
	row := y1
	col := x1
	for i, r := range []rune(typeTest.Text) {
    if i >= typeTest.CurPos {
      s.SetContent(col, row, r, nil, style)
    }else if typeTest.Results[i] {
      s.SetContent(col, row, r, nil, style.Foreground(tcell.ColorGreen))
    }else {
      s.SetContent(col, row, r, nil, style.Foreground(tcell.ColorRed))
    }

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

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) {
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

}

