package parser

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"golang.org/x/net/html"
	"html/template"
	"regexp"
	"log"
)

type UserPosts struct {
	url string
	minKarma int
}

func (p *UserPosts) GetPagination() (currentPage, pageMax int, err error) {
	doc, err := goquery.NewDocument(p.url + "/content/&type=forums_topic_post")
	if err != nil {
		return
	}
	log.Print(p.url)

	pSelect := doc.Find(".ipsPagination_pageJump a")
	if pSelect.Length() > 0 {
		pSelect.Each(func(i int, s *goquery.Selection) {
			paginatorText := s.Text()
			regular, localErr := regexp.Compile("Страница (?P<Page>\\d+) из (?P<MaxPage>\\d+)")
			if localErr != nil {
				err = localErr
			}

			result := regular.FindStringSubmatch(paginatorText)
			if len(result) > 2 {
				currentPage, _ = strconv.Atoi(result[1])
				pageMax, _ = strconv.Atoi(result[2])
			} else {
				err = errors.New("regxp pagination not found")
			}
		})
	} else {
		err = errors.New("pagination not found")
	}
	return
}

func (p *UserPosts) GetPage(pageNum int) ([]*Message, error) {
	doc, err := goquery.NewDocument(p.url + "/content/&type=forums_topic_post&page=" + strconv.Itoa(pageNum))
	if err != nil {
		return nil, err
	}

	aSelect := doc.Find("article.ipsComment")

	if aSelect.Length() > 0 {
		var messages []*Message

		author := strings.TrimSpace(doc.Find("#elProfileHeader > div.ipsColumns.ipsColumns_collapsePhone > div.ipsColumn.ipsColumn_fluid > div > h1").Text())
		aSelect.Each(func(i int, s *goquery.Selection) {
			subS := s.Find("div[data-role='reactCount']")

			karma, localErr := strconv.Atoi(strings.TrimSpace(subS.Text()))
			if localErr != nil {
				err = localErr
			}

			if karma >= p.minKarma {
				s.Find("div .ipsPad_half > .ipsPad > div[data-role='commentContent'] blockquote").Each(func(i int, selection *goquery.Selection) {
					selection.Nodes[0].Attr = []html.Attribute{}
					//selection.SetText(">> " + selection.Text())
				})

				messageText, localErr := s.Find("div .ipsPad_half > .ipsPad > div[data-role='commentContent']").Html()
				if localErr != nil {
					err = localErr
					return
				}
				time, _ := s.Find("div.ipsPad_half > div > p > a > time").Attr("datetime")
				message := &Message{}
				message.Karma = karma
				message.Body = template.HTML(messageText)
				message.Author = author
				message.Date = time

				messages = append(messages, message)
			}
		})

		if err != nil {
			return nil, err
		}

		return messages, nil
	} else {
		return nil, errors.New("article not found")
	}
}

func NewUserPostsParser(url string, minKarma int) ParserInterface {
	return &UserPosts{
		url: url,
		minKarma: minKarma,
	}
}
