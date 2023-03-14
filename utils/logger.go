package utils

import (
	"fmt"
	"log"
	"time"
)

func TimeLog(message string) {
	t := time.Now()
	formattedTime := t.Format("02.01.2006 15:04")

	fmt.Printf("%s# %s \n", formattedTime, message)
	log.Printf("%s# %s \n", formattedTime, message)
}
