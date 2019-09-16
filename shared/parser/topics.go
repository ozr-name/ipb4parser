package parser

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
	"strings"
	"html/template"
)

type TopicParser struct {
	url string
	minKarma int
}

func (p *TopicParser) GetPagination() (currentPage, pageMax int, err error) {
	doc, err := goquery.NewDocument(p.url)
	if err != nil {
		return
	}

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

func (p *TopicParser) GetPage(pageNum int) ([]*Message, error) {
	doc, err := goquery.NewDocument(p.url + "&page=" + strconv.Itoa(pageNum))
	if err != nil {
		return nil, err
	}

	aSelect := doc.Find("article.ipsComment")

	if aSelect.Length() > 0 {
		var messages []*Message

		aSelect.Each(func(i int, s *goquery.Selection) {
			subS := s.Find("div[data-role='reactCount']")

			karma, localErr := strconv.Atoi(strings.TrimSpace(subS.Text()))
			if localErr != nil {
				err = localErr
			}

			if karma > p.minKarma {

				messageText, localErr := s.Find("div > div.cPost_contentWrap.ipsPad > div.ipsType_normal.ipsType_richText.ipsContained").Html()
				if localErr != nil {
					err = localErr
					return
				}

				author := s.Find("aside > h3 > strong > a").Text()
				time, _ := s.Find("div > div.ipsComment_meta.ipsType_light > div:nth-child(2) > a > time").Attr("datetime")

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

func NewTopicParser(url string, minKarma int) ParserInterface {
	return &TopicParser{
		url: url,
		minKarma: minKarma,
	}
}
