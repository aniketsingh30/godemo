package logger

import "log"

func PANIC(message string, err error) {
	if err != nil {
		log.Panic(message, err)
	}
}

func INFO(message string, data interface{}) {
	log.Println(message, data)
}

func INFOMessage(message string) {
	log.Println(message)
}


func ErrorMessage(message string) {
	log.Println(message)
}