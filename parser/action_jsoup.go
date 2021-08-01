package parser

/**
Jsoup解析
*/

type JsoupParser struct {
	Action
}

func (parser JsoupParser) parseEach(rule string, needFilterString bool) []string {
	return []string{"hello", "world"}
}
