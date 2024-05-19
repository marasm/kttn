package typingTest

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
	"github.com/marasm/kttn/internal/utils"
)


type TypingTest struct {
  Text string
  Results []bool 
  CurPos int
  StartTime time.Time
  EndTime time.Time
}

func CreateNewTest(numOfWords int, wordSet string)  TypingTest {
  text := utils.GetWordsFromFile(numOfWords, wordSet)
  return TypingTest{
    Text: text,
    CurPos: 0,
    Results: make([]bool, utf8.RuneCountInString(text)),
  } 
}

func (t TypingTest) TestComplete() bool {
  return t.CurPos >= t.GetTotalChars() - 1  
}

func (t *TypingTest) UpdateWithRegKey(key rune) {
  // - see if the key == the current rune in text
  // - append the results with true|false
  t.Results[t.CurPos] = []rune(t.Text)[t.CurPos] == key
  if !t.TestComplete() {
    t.CurPos++
  }
  if t.StartTime.IsZero()  && t.CurPos > 0 {
    t.StartTime = time.Now()
  }
  if t.EndTime.IsZero()  && t.TestComplete() {
    t.EndTime = time.Now()
  }
}

func (t *TypingTest) UpdateWithBackspace(key tcell.Key) {
  if t.CurPos > 0 {
    t.CurPos--
  }
}

func (t TypingTest) GetErrorCount() int {
  c := 0 
  for _, r := range t.Results[:t.CurPos + 1] {
    if !r {
      c++
    }
  }
  return c
}

func (t TypingTest) GetTotalChars() int {
  return utf8.RuneCountInString(t.Text)
}

func (t TypingTest) GetWordCount() int {
  return len(strings.Split(t.Text, " "))
}

func (t TypingTest) GetAccuracyPercent() float32 {
  total := t.GetTotalChars() 
  success := total - t.GetErrorCount()
  return float32(success)/float32(total)*100
}

func (t TypingTest) GetResultsAsString() string {
  cpm := float64(t.GetTotalChars())/t.EndTime.Sub(t.StartTime).Minutes()
  wpm := float64(t.GetWordCount())/t.EndTime.Sub(t.StartTime).Minutes()
  total := t.GetTotalChars()
  errors := t.GetErrorCount()
  accuracy := t.GetAccuracyPercent()
  return fmt.Sprintf(`
         WPM : %.2f
         CPM : %.2f
 Total Typed : %d
      Errors : %d
    Accuracy : %.2f%%`, 
    wpm, cpm, total, errors, accuracy)
}

