package parser

import "html/template"

type ParserInterface interface {
	GetPagination() (currentPage, pageMax int, err error)
	GetPage(pageNum int) ([]*Message, error)
}

type Message struct {
	Karma int
	Body template.HTML
	Author string
	Date string
}