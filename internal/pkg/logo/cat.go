package logo

import (
	"math/rand/v2"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/marasm/kttn/internal/pkg/utils"
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

func UpdateLogo(screen tcell.Screen, style tcell.Style) {
  x, y := utils.GetLogoCoords(screen)
  n := rand.IntN(5)
  switch n {
    case 0:
      utils.DrawText(screen, x, y, style, LOGO)
    case 1:
      utils.DrawText(screen, x, y, style, GetLogoWithParams(WATCHING_EYE, TAIL_RIGHT_FLAT))
    case 2:
      utils.DrawText(screen, x, y, style, GetLogoWithParams(HAPPY_EYE, TAIL_LEFT_FLAT))
    case 3:
      utils.DrawText(screen, x, y, style, GetLogoWithParams(SURPRISED_EYE, TAIL_RIGHT_UP))
    case 4:
      utils.DrawText(screen, x, y, style, GetLogoWithParams(SLEEPY_EYE, TAIL_LEFT_UP))
    
  }
}

func GetLogoWithParams(eyes, tail string) string {
  r := strings.Replace(LOGO, "^", eyes, 2)
  return strings.Replace(r, TAIL_BEHIND, tail, 1)
}


