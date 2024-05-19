package utils

import (
	"bufio"
	"math/rand/v2"
	"os"
	"strings"

	_ "github.com/marasm/kttn/internal/init"
)


func GetWordsFromFile(numOfWords int, fileName string) string {
  var res strings.Builder
  words := readAllWordsFromFile(fileName)
  for i := 0; i < numOfWords; i++ {
    //add a space between words
    if i > 0 {
      res.WriteString(" ")
    }
    res.WriteString(words[rand.IntN(len(words))])
  }
  return res.String()
}

func readAllWordsFromFile(fileName string) []string {
  var results []string

  file, err := os.Open("words/" + fileName)
  if err != nil {
    println("Error opening file" + err.Error())
  }

  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)
  
  for scanner.Scan() {
    word := strings.TrimSpace(scanner.Text())
    if len(word) > 0 {
      results = append(results,word) 
    }
  }

  return results

}

