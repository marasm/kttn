package utils

func DistanceToNextWhitespace(curPos int, textAsRunes []rune) int {
  res := 0
  for i, r := range textAsRunes[curPos + 1:] {
    if r == ' ' || curPos + i >= len(textAsRunes) - 2 {
      res = i + 1
      break
    }
  }
  return res
}
