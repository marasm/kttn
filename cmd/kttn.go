package main

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/gdamore/tcell/v2"
	tt "github.com/marasm/kttn/internal/pkg/typingTest"
	utils "github.com/marasm/kttn/internal/pkg/utils"
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


type Cat struct {
  CurEyes string
  CurTail string
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

  typeTest := tt.CreateNewTest()

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
      showLegend(s, defStyle)
      updateTypingBox(s, defStyle, typeTest)
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape  {
        typeTest = tt.CreateNewTest()
        s.Clear()
        updateLogo(s, defStyle)
        showLegend(s, defStyle)
        updateTypingBox(s, defStyle, typeTest)
      } else if ev.Key() == tcell.KeyCtrlW || ev.Key() == tcell.KeyCtrlQ {
				return
			} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
        typeTest.UpdateWithBackspace(ev.Key())
        updateTypingBox(s, defStyle, typeTest)
			} else {
        if typeTest.TestComplete() {
          typeTest.UpdateWithRegKey(ev.Rune())
          s.Clear()
          updateLogo(s, defStyle)
          showResults(s, defStyle, typeTest)
        } else {
          typeTest.UpdateWithRegKey(ev.Rune())
          updateLogo(s, defStyle)
          updateTypingBox(s, defStyle, typeTest)
        }
      }
		}
	}
}

func updateLogo(screen tcell.Screen, style tcell.Style) {
  x, y := utils.GetLogoCoords(screen)
  n := rand.IntN(5)
  switch n {
    case 0:
      utils.DrawText(screen, x, y, style, LOGO)
    case 1:
      utils.DrawText(screen, x, y, style, getLogoWithParams(WATCHING_EYE, TAIL_RIGHT_FLAT))
    case 2:
      utils.DrawText(screen, x, y, style, getLogoWithParams(HAPPY_EYE, TAIL_LEFT_FLAT))
    case 3:
      utils.DrawText(screen, x, y, style, getLogoWithParams(SURPRISED_EYE, TAIL_RIGHT_UP))
    case 4:
      utils.DrawText(screen, x, y, style, getLogoWithParams(SLEEPY_EYE, TAIL_LEFT_UP))
    
  }
}

func getLogoWithParams(eyes, tail string) string {
  r := strings.Replace(LOGO, "^", eyes, 2)
  return strings.Replace(r, TAIL_BEHIND, tail, 1)
}

func updateTypingBox(screen tcell.Screen, style tcell.Style, typeTest tt.TypingTest) {
  sx, sy, ex,ey := utils.GetTypingBoxCoords(screen, typeTest.Text) 
  utils.DrawBox(screen, sx, sy, ex, ey, style)
	utils.DrawBoundedText(screen, sx+1, sy+1, ex, ey-1, style, typeTest)
}

func showLegend(screen tcell.Screen, style tcell.Style) {
  midX, _ := utils.GetMidScreenCoords(screen)
  _, maxY := screen.Size()
  utils.DrawText(screen, midX - 20, maxY - 1, style, "C-q or C-w to quit | Esc to restart the test")
}

func showResults(screen tcell.Screen, style tcell.Style, typeTest tt.TypingTest) {
  screen.HideCursor()  
  midX, midY := utils.GetMidScreenCoords(screen)
  cpm := float64(typeTest.GetTotalChars())/typeTest.EndTime.Sub(typeTest.StartTime).Minutes()
  wpm := float64(typeTest.GetWordCount())/typeTest.EndTime.Sub(typeTest.StartTime).Minutes()
  total := typeTest.GetTotalChars()
  errors := typeTest.GetErrorCount()
  accuracy := typeTest.GetAccuracyPercent()
  resultsStr := fmt.Sprintf(`
         WPM : %.2f
         CPM : %.2f
 Total Typed : %d
      Errors : %d
    Accuracy : %.2f%%`, 
    wpm, cpm, total, errors, accuracy)
  utils.DrawText(screen, midX - 14, midY - 2, style, resultsStr)
}

