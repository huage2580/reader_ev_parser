package ev_test

import (
	"fmt"
	"reader_ev_parser/parser"
	"testing"
)

func Test_ME(t *testing.T) {
	//Jsoup()
	//m01xs()
	//m01xcList()
	//Biquge()
	KenShuWu()
}

func Jsoup() {
	//jsoupListStr("tag.h2.0@ownText")
	//jsoupListStr("tag.h1.0@ownText")
	//jsoupListStr("text.甲方@html")
	//jsoupListStr("id.list.0@tag.dd.0:1:2@text&&id.list.0@tag.dd.12:13@text&&id.list.0@tag.dd.12:13@html")
	//jsoupListStr("tag.h1.0@html")
	//jsoupListStr("id.list.0@tag.dd.0@all")
	//jsoupListStr("id.list.0@tag.dd.0@html")
	//jsoupListStr("id.list.0@tag.dd.0@tag.a.0@href")

	//jsoupBatch("id.list.0@tag.dd.0:1:2:3:4", "tag.a@href", "text")
	//jsoupHtml("tag.div@html")
	//jsoupHtml("##div##fxxk###")
	//jsoupBatch("id.list.0@tag.dd.0:1:2:3:4", "children@tag", "tag")
	//KenShuWu()
	Biquge()
}

func jsoupListStr(rule string) {
	tId := parser.StartTransaction(DATA)
	result := parser.ParseRuleStr(tId, rule)
	fmt.Printf("jsoup rule-> [%s] result-> %s\n", rule, result)
	parser.EndTransaction(tId)
}

func jsoupBatch(rule string, rule1 string, rule2 string) {
	tId := parser.StartTransaction(DATA)
	bId := parser.ParseRuleRaw(tId, rule)
	size := parser.QueryBatchResultSize(bId)
	fmt.Printf("jsoup rule-> [%s] result-> %d\n", rule, size)
	for i := 0; i < size; i++ {
		r1 := parser.ParseRuleStrForParent(bId, rule1, i)
		fmt.Printf("jsoup rule-> [%s] result-> %s\n", rule1, r1)
		r2 := parser.ParseRuleStrForParent(bId, rule2, i)
		fmt.Printf("jsoup rule-> [%s] result-> %s\n", rule2, r2)
	}
	parser.EndTransaction(bId)
	parser.EndTransaction(tId)
}

func jsoupHtml(rule string) {
	tId := parser.StartTransaction(`
<div><p>hello</p><p>wor<br/>ld</p>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;张若尘现在所在的国家，名叫“云武郡国”，只是昆仑界东域成千上万个郡国中的一个。<br>测试</div>
`)
	result := parser.ParseRuleStr(tId, rule)
	fmt.Printf("jsoup rule-> [%s] result-> %s\n", rule, result)
	parser.EndTransaction(tId)
}

func Regexp() {
	var test = parser.RegexpFilter([]string{"测试内容,测试内容，测试"}, "##测试")
	fmt.Println(test)
	var test2 = parser.RegexpFilter([]string{"测试内容,测屁内容，测试"}, "##测(.)##体$1")
	fmt.Println(test2)
	var test3 = parser.RegexpFilter([]string{"测试内容,测试内容，测试"}, "##(.)试##替换$1的###")
	fmt.Println(test3)
}

func KenShuWu() {
	//搜索
	jsoupBatchInput(DATA_KENSHUWU_SEARCH, "class.novelslist2@li!0", []string{"tag.a.0@text", "tag.span.2@text", "a@href", "tag.a.1@text"})
	//章节
	jsoupBatchInput(DATA_KENSHUWU_LIST, "id.list@dd", []string{"a@text", "a@href"})

	//阅读页
	jsoupStrInput(DATA_KENSHUWU_DETAIL, "id.content@textNodes")

}

func Biquge() {
	//搜索
	//jsoupBatchInput(DATA_KENSHUWU_SEARCH,"class.novelslist2@li!0",[]string{"tag.a.0@text","tag.span.2@text","a@href","tag.a.1@text"})
	//章节
	//jsoupBatchInput(DATA_KENSHUWU_LIST,"id.list@dd",[]string{"a@text","a@href"})
	//jsoupBatchInput(DATA_KENSHUWU_LIST, "id.list@dd", []string{"a@text", "a@href"})
	jsoupBatchInput(DATA_biquge_LIST, ".mulu_list li a", []string{"text", "href"})

	//阅读页
	//jsoupStrInput(DATA_KENSHUWU_DETAIL, "id.content@textNodes")

}

func m01xs() {
	jsoupBatchInput(DATA_01xs_SEARCH, ".search-list dl",
		[]string{"tag.a.0@text", "tag.dd.0@text##.*：|\\s.*", "tag.dd.0:1@text##.*\\s", "tag.a.1@text", "tag.a.0@href", "tag.a.0@href##.+\\D((\\d+)\\d{3})\\D##https://img01xs.cdn.bcebos.com/files/article/image/$2/$1/$1s.jpg###"})

}

func m01xcList() {
	jsoupBatchInput(DATA_m01xs_LIST, "#index_list!-1@li a", []string{"text", "href"})

}

func jsoupBatchInput(input string, rule string, rule1 []string) {
	tId := parser.StartTransaction(input)
	bId := parser.ParseRuleRaw(tId, rule)
	size := parser.QueryBatchResultSize(bId)
	fmt.Printf("jsoup rule-> [%s] result-> %d\n", rule, size)
	for i := 0; i < size; i++ {
		for _, r := range rule1 {
			r1 := parser.ParseRuleStrForParent(bId, r, i)
			fmt.Printf("%s", r1)
		}
		fmt.Println("")
	}
	parser.EndTransaction(bId)
	parser.EndTransaction(tId)
}

func jsoupStrInput(input string, rule string) {
	tId := parser.StartTransaction(input)
	result := parser.ParseRuleStr(tId, rule)
	fmt.Printf("jsoup rule-> [%s] result-> \n", rule)
	for i, s := range result {
		fmt.Printf("%d->%s\n", i, s)
	}
	parser.EndTransaction(tId)
}
