package main

import (
	"flag"
	"gopkg.in/cheggaaa/pb.v1"
	"html/template"
	"log"
	"os"
//	"github.com/ruelephant/ipb4parser/shared/parser"
	"github.com/ozr-name/ipb4parser/shared/parser"
)

// Flags Variables
var topic string
var user string
var karmaMin int

func init() {
	flag.StringVar(&topic, "topic", "", "Введите url топика, который хотите спарсить")
	flag.StringVar(&user, "user", "", "Введите url пользователя, которого хотите спарсить")
	flag.IntVar(&karmaMin, "count", 0, "Количество лайков, от которых собирать. По умолчанию значение = 0. Если поставить 2, то собираться будет от 3-ёх.")
	flag.Parse()
}

func main() {
	var p parser.ParserInterface // Преинициализация, если мы обьявим p внутри {} мы не получим ее за её пределами

	// Двойное условие для избежания не однозначной ситуации когда заполнены оба параметра
	if topic != "" && user == "" {
		p = parser.NewTopicParser(topic, karmaMin)
	} else if topic == "" && user != "" {
		p = parser.NewUserPostsParser(user, karmaMin)
	} else {
		log.Fatal("Введите url топика (--topic) или пользователя (--user), которые Вы хотите спарсить!")
	}

	_, maxPage, err := p.GetPagination()
	if err != nil {
		log.Fatal(err)
	}

	var messages []*parser.Message

	bar := pb.StartNew(maxPage)

	for i := 1; i <= maxPage; i++ {
		bar.Increment()
		messagesInPage, err := p.GetPage(i)
		if err != nil {
			log.Fatal(err)
		}

		messages = append(messages, messagesInPage...)
	}
	bar.Finish()

	t, err := template.ParseFiles("resources/tmpl/list.html")
	if err != nil {
		log.Fatal("Parse error: ", err)
	}

	f, err := os.Create("result.html")
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	err = t.Execute(f, messages)
	if err != nil {
		log.Print("Execute error: ", err)
	}
}
