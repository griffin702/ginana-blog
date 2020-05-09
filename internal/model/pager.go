package model

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Pager struct {
	Page      int64             `json:"page"`
	PageSize  int64             `json:"page_size"`
	AllPage   int64             `json:"all_page"`
	AllCount  int64             `json:"all_count"`
	UrlPath   string            `json:"url_path"`
	UrlParams map[string]string `json:"url_params"`
}

func (p *Pager) SetArticleID(id int64, paramName ...string) {
	if id > 0 {
		if len(paramName) > 0 { // 放置到UrlParams里
			p.UrlParams[paramName[0]] = strconv.FormatInt(id, 10)
		} else { // 放置到UrlPath里
			p.UrlPath = fmt.Sprintf("%s/%d", p.UrlPath, id)
		}
	}
}

func (p *Pager) url(page int64) string {
	p.UrlParams["page"] = strconv.FormatInt(page, 10)
	p.UrlParams["pagesize"] = strconv.FormatInt(p.PageSize, 10)
	var params []string
	for k, v := range p.UrlParams {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}
	paramsStr := strings.Join(params, "&")
	return fmt.Sprintf("%s?%s", p.UrlPath, paramsStr)
}

func (p *Pager) ToString() string {
	var buf bytes.Buffer
	var from, to, limitLink, offset, totalpage int64
	var omit string
	offset = 2
	limitLink = 4
	if p.Page < 3 {
		limitLink = 5
	}
	totalpage = int64(math.Ceil(float64(p.AllCount) / float64(p.PageSize)))
	p.AllPage = totalpage
	if totalpage < limitLink {
		from = 1
		to = totalpage
	} else {
		from = p.Page - offset
		to = from + limitLink
		if from < 1 {
			from = 1
			to = from + limitLink - 1
		} else if to > totalpage {
			to = totalpage
			from = totalpage - limitLink + 1
		}
	}
	buf.WriteString("<ul class=\"pagination\">")
	if p.Page > 1 {
		buf.WriteString(fmt.Sprintf("<li><a href=\"%s\">上一页</a></li>", p.url(p.Page-1)))
	} else {
		buf.WriteString("<li class=\"disabled\"><a>上一页</a></li>")
	}
	if p.Page >= limitLink {
		if p.Page-limitLink > 0 && totalpage != 5 {
			omit = "..."
		}
		if totalpage != 4 {
			buf.WriteString(fmt.Sprintf("<li><a href=\"%s\">1%s</a></li>", p.url(1), omit))
		}
	}
	for i := from; i <= to; i++ {
		if i == p.Page {
			buf.WriteString(fmt.Sprintf("<li class=\"active\"><a>%d</a></li>", i))
		} else {
			buf.WriteString(fmt.Sprintf("<li><a href=\"%s\">%d</a></li>", p.url(i), i))
		}
	}
	if totalpage > to {
		if totalpage-to > 1 {
			omit = "..."
		}
		buf.WriteString(fmt.Sprintf("<li><a href=\"%s\">%s%d</a></li>", p.url(totalpage), omit, totalpage))
	}
	if p.Page < totalpage {
		buf.WriteString(fmt.Sprintf("<li><a href=\"%s\">下一页</a></li>", p.url(p.Page+1)))
	} else {
		buf.WriteString(fmt.Sprintf("<li class=\"disabled\"><a>下一页</a></li>"))
	}
	buf.WriteString("</ul>")
	return buf.String()
}
