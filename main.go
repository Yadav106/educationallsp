package main

import (
	"bufio"
	"log"
	"os"

	"github.com/Yadav106/educationallsp/rpc"
)

func main() {
  logger := getLogger("/Users/macbook/Desktop/Programming/educationallsp/log.txt")
  logger.Println("Mic Check! 1! 2! 3!")
  scanner := bufio.NewScanner(os.Stdin)
  scanner.Split(rpc.Split)

  for scanner.Scan() {
    msg := scanner.Text()
    handleMessage(logger, msg)
  }
}

func handleMessage (logger *log.Logger, msg any) {
  logger.Println(msg)
}

func getLogger(fileName string) *log.Logger {
  logfile, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
  if err != nil {
    panic("give a better file 🗿")
  }

  return log.New(logfile, "[educationallsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
