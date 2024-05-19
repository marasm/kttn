package utils

import (
  "strings"
  "testing"
)


func TestGetWordsFromFileSuccess(t *testing.T) {
  res := GetWordsFromFile(10, "en_500")
  wordCount := strings.Split(res, " ")
  if len(wordCount) != 10 {
    t.Fatalf("Expected 10 words but got %d", len(wordCount))
  }
}
