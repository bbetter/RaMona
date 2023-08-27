package utils

import (
	"fmt"
	"log"
	"time"
)

func TimeLog(message string) {
	t := time.Now()
	// dates in GoLang are bOOlshit
	formattedTime := t.Format("02.01.2006 15:04")

	fmt.Printf("%s# %s \n", formattedTime, message)
	log.Printf("%s# %s \n", formattedTime, message)
}
