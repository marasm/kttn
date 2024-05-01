package utils

import (
	"bufio"
	"os"
	"strings"
)


func GetWordsFromFile(numOfWords int, fileName string) string {
  cwd, _ := os.Getwd() 
  file, err := os.Open(cwd + "/words/" + fileName)
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
    res.WriteString(scanner.Text())
  }
  return res.String()
}
