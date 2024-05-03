package screen

import (
	"math"
	"strings"
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
	tt "github.com/marasm/kttn/internal/pkg/typingTest"
	"github.com/marasm/kttn/internal/pkg/utils"
) 


func GetMidScreenCoords(screen tcell.Screen) (midX int, midY int) {
  availX, availY := screen.Size()
  return availX/2, availY/2
}

func GetLogoCoords(screen tcell.Screen) (logoX, logoY int) {
  midX, midY := GetMidScreenCoords(screen)
  return midX - 6, midY - 8 
}

func GetTypingBoxCoords(screen tcell.Screen, text string) (startX, startY, endX, endY int) {
  maxX, _ := screen.Size()
  _, midY := GetMidScreenCoords(screen) 
  txtLength := utf8.RuneCountInString(text)
  //add .75 to always round up and to account for word wrapping line breaks
  numOfRows := math.Round(float64(txtLength)/float64(maxX-12) + 0.75)  
  offset := 0
  if numOfRows > 6 {
    offset = int(numOfRows) - 6
  }
  return 5, midY - int(math.Round(numOfRows/2 + 0.75)) + offset, maxX - 5, midY + int(math.Round(numOfRows/2 + 0.75)) + offset
}

func UpdateTypingBox(screen tcell.Screen, style tcell.Style, typeTest tt.TypingTest) {
  sx, sy, ex,ey := GetTypingBoxCoords(screen, typeTest.Text) 
  DrawBox(screen, sx, sy, ex, ey, style)
	DrawBoundedText(screen, sx+1, sy+1, ex, ey-1, style, typeTest)
}

func DrawLegend(screen tcell.Screen, style tcell.Style) {
  midX, _ := GetMidScreenCoords(screen)
  _, maxY := screen.Size()
  DrawText(screen, midX - 20, maxY - 1, style, "C-q or C-w to quit | Esc to restart the test")
}

func DrawResults(screen tcell.Screen, style tcell.Style, typeTest tt.TypingTest) {
  screen.HideCursor()  
  midX, midY := GetMidScreenCoords(screen)
  DrawText(screen, midX - 14, midY - 2, style, typeTest.GetResultsAsString())
}

func DrawText(s tcell.Screen, x, y int, style tcell.Style, text string) {
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

func DrawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) {
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

func DrawBoundedText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, typeTest tt.TypingTest) {
	row := y1
	col := x1
  textAsRunes := []rune(typeTest.Text)
	for i, r := range textAsRunes {
    if i >= typeTest.CurPos {
      s.SetContent(col, row, r, nil, style.Foreground(tcell.ColorGray))
    }else if typeTest.Results[i] {
      s.SetContent(col, row, r, nil, style)
    }else {
      s.SetContent(col, row, r, nil, style.Foreground(tcell.ColorRed))
    }

    if i == typeTest.CurPos {
      s.ShowCursor(col, row)
    }

		col++
    // - if cur column about to cross the right border (x2) == break
    // - if currently on whitespace and the distance to next whitespece > distance to x2 == break
		if col >= x2 || 
      (r == ' ' && col + utils.DistanceToNextWhitespace(i, textAsRunes) >= x2) {
			row++
			col = x1
		}

		if row > y2 {
			break
		}
	}
}
