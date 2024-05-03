package utils

import (
	"bufio"
	"os"
	"strings"

  _ "github.com/marasm/kttn/internal/pkg/init"
)


func GetWordsFromFile(numOfWords int, fileName string) string {
  file, err := os.Open("words/" + fileName)
  if err != nil {
    println("Error opening file" + err.Error())
  }
  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)
  var res strings.Builder
  //TODO add randomness to this
  for i := 0; i < numOfWords; i++ {
    success := scanner.Scan()
    if success == false {
      break
    }
    //add a space between words
    if i > 0 {
      res.WriteString(" ")
    }
    res.WriteString(scanner.Text())
  }
  return res.String()
}
