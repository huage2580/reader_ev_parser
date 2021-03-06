package parser

const PARSER_TYPE_CSS = "^[-+]?@css:"
const PARSER_TYPE_JSON = `^[-+]?@JSon:|^[-+]?@json:|^\$\.`
const PARSER_TYPE_XPATH = `^[-+]?@XPath:|^//`
const PARSER_TYPE_REGEXP = "^[-+]?:" //正则之AllInOne

const RE_REPLACE = "##"
const RE_JS_TOKEN = `@js:(.*)|<js>(.*)<\/js>`

const OPERATOR_AND = "&&"   //组合拼接
const OPERATOR_OR = "||"    //或
const OPERATOR_MERGE = "%%" //依次，左一个，右一个

const DELIMITER = "@"

const FILTER_TEXT = "text"
const FILTER_TEXT_NODE = "textNodes"
const FILTER_OWN_TEXT = "ownText"
const FILTER_HTML = "html"
const FILTER_ALL = "all"
const FILTER_HREF = "href"
const FILTER_TAG_NAME = "tag" //only for debug

const JSOUP_SPLIT = "."

const JSOUP_SUPPORT_CHILD = "children"
const JSOUP_SUPPORT_CLASS = "class"
const JSOUP_SUPPORT_TAG = "tag"
const JSOUP_SUPPORT_ID = "id"
const JSOUP_SUPPORT_TEXT = "text"

const JSOUP_EXCLUDE_CHAR = "!"
const JSOUP_EXCLUDE_INT = ":"

const ACTION_TYPE_JSOUP = "jsoup"
const ACTION_TYPE_CSS = "css"

const REGEXP_FOR_JSOUP_1 = `.*\.[-]?\d` //有些奇怪的只有两段的jsoup格式，例如 a.1,匹配之后执行前面的css再筛选后面的
