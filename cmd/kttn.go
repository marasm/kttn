package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"

	config "github.com/marasm/kttn/internal/config"
	logo "github.com/marasm/kttn/internal/logo"
	scr "github.com/marasm/kttn/internal/screen"
	tt "github.com/marasm/kttn/internal/typingTest"
)


func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Printf("%+v", err)
	}
	if err := s.Init(); err != nil {
		fmt.Printf("%+v", err)
	}

  conf := config.InitConfig()

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

  typeTest := tt.CreateNewTest(conf.NumberOfWords, conf.WordSet)

	for {
    // Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
      s.Clear()
      logo.UpdateLogo(s, defStyle)
      scr.DrawLegend(s, defStyle)
      scr.UpdateTypingBox(s, defStyle, typeTest)
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape  {
        typeTest = tt.CreateNewTest(conf.NumberOfWords, conf.WordSet)
        s.Clear()
        logo.UpdateLogo(s, defStyle)
        scr.DrawLegend(s, defStyle)
        scr.UpdateTypingBox(s, defStyle, typeTest)
      } else if ev.Key() == tcell.KeyCtrlW || ev.Key() == tcell.KeyCtrlQ {
				return
			} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
        typeTest.UpdateWithBackspace(ev.Key())
        scr.UpdateTypingBox(s, defStyle, typeTest)
			} else {
        if typeTest.TestComplete() {
          typeTest.UpdateWithRegKey(ev.Rune())
          s.Clear()
          logo.UpdateLogo(s, defStyle)
          scr.DrawResults(s, defStyle, typeTest)
          if typeTest.GetWpm() > conf.MaxWpm {
            conf.MaxWpm = typeTest.GetWpm()
          }
          if typeTest.GetCpm() > conf.MaxCpm {
            conf.MaxCpm = typeTest.GetCpm()
          }
          config.SaveConfig(conf)
        } else {
          typeTest.UpdateWithRegKey(ev.Rune())
          logo.UpdateLogo(s, defStyle)
          scr.UpdateTypingBox(s, defStyle, typeTest)
        }
      }
		}
	}
}


