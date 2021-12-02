package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

package main

import (
"fmt"
"golang.org/x/net/html"
"io/ioutil"
"net/http"
"net/url"
"strings"
)

// DeadLink 题目要求构造的死链结构，由死链url和上一级url组成
type DeadLink struct {
	DeadLink string
	SourceLink string
}

//var ans []DeadLink// 最后返回的答案

var DeadLinks map[string]DeadLink

func DeadLinkDetector(rootUrl string) []DeadLink {
	// your implementation
	DeadLinks = make(map[string]DeadLink)
	DeadLinks = help(rootUrl,DeadLinks)
	//ans := make([]DeadLink, 0)
	var ans []DeadLink
	for _, link := range DeadLinks {
		//temp := link.
		fmt.Printf("DeadLink: %s\r\n",link.DeadLink)
		fmt.Printf("SourceLink: %s\r\n",link.SourceLink)
		ans = append(ans,link)
	}
	return ans
}

// GetResCode 获取url网址的响应码resCode
func GetResCode(rootUrl string) int{
	u, _ := url.Parse(rootUrl)//将string解析成url格式
	q := u.Query()//将路径解析为一个方便操作的对象
	u.RawQuery = q.Encode()//将处理后url的回传给u
	res, err := http.Get(u.String())
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return 0
	}
	resCode := res.StatusCode
	return resCode
}

// GetChildLink 解析得到url网址的*html.Node，方便visit()访问遍历
func GetChildLink(rootUrl string) *html.Node{
	u, _ := url.Parse(rootUrl)//将string解析成url格式
	q := u.Query()//将路径解析为一个方便操作的对象
	u.RawQuery = q.Encode()//将处理后url的回传给u
	res, err := http.Get(u.String())
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)//取出响应的body
	if err != nil {
		fmt.Printf("read from res.Body failed, err:%v\n", err)
		return nil
	}
	b := string(body) // 转换成string类型
	doc, err := html.Parse(strings.NewReader(b))
	if err != nil {
		fmt.Printf("findlinks: %v\n", err)
		return nil
	}
	return doc
}

// visit 对*html.Node的内容迭代+遍历，得到该url下的网址列表
func visit(links map[string]interface{}, n *html.Node) map[string]interface{} {
	// 用links这个map存储当前html下的网址
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" && strings.HasPrefix(a.Val, "http"){
				//resCode := GetResCode(a.Val)
				links[a.Val] = struct{}{}
			}
		}
	}
	//这个递归，是对当前报文从上到下从里到外、逐级逐行检测
	//FirstChild是第一个子节点，NextSibling是下一个同级节点
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

// help 调用GetChildLink和visit，实现判断输入url下面两级死链的功能
func help(rootUrl string, DeadLinks map[string]DeadLink)  map[string]DeadLink{
	//DeadLinks := make(map[string]DeadLink)// map类型的DeadLinks无重复地储存所有DeadLink
	tempLinks1 := make(map[string]interface{})// map类型的tempLinks1暂存html中的所有链接
	doc := GetChildLink(rootUrl)// 解析得到url网址的*html.Node
	linksLayer1 := visit(tempLinks1, doc)// 得到第一层所有的链接列表，用map类型的linkLayer1储存
	fmt.Printf("当前网页有%d个一级链接\r\n", len(linksLayer1))
	// 用link取出第一层链接，挨个判断是否为死链
	for link := range linksLayer1 {//使用length?
		//go func() {
		fmt.Printf("一级link： %s\r\n  ", link)
		resCode := GetResCode(link)
		fmt.Printf("%d\r\n", resCode)
		if resCode < 200 || resCode >= 300 { // 情况1：判断为死链
			//fmt.Printf("%s ", link)
			//fmt.Printf("%d\r\n", resCode)
			// 借助temp，把死链加入DeadLinks
			var temp DeadLink
			temp.DeadLink = link
			temp.SourceLink = rootUrl
			DeadLinks[link] = temp

		} else { // 情况2：不是死链，继续获得第二层链接，并判断是否为死链
			//tempLinks2 := make(map[string]interface{}) //一点疑问，好像不合适，是会每次重新声明吗？
			//doc := GetChildLink(link)
			//linksLayer2 := visit(tempLinks2, doc) // 得到该link对应的第二层链接列表，用map类型的linkLayer2储存
			//fmt.Printf("当前网页有%d个二级链接\r\n", len(linksLayer2))
			//// 用link取出第二层链接，挨个判断是否为死链
			//for link := range linksLayer2 {
			//	fmt.Printf("二级link： %s\r\n", link)
			//	resCode := GetResCode(link)
			//	fmt.Printf("%d\r\n", resCode)
			//	if resCode < 200 || resCode >= 300 {
			//		//fmt.Printf("%s ", link)
			//		//fmt.Printf("%d\r\n", resCode)
			//		// 借助temp，把死链加入DeadLinks
			//		var temp DeadLink
			//		temp.DeadLink = link
			//		temp.SourceLink = rootUrl
			//		DeadLinks[link] = temp
			//	}
			//}
			help(link, DeadLinks)
		}
		//}()
		//time.Sleep(1 * time.Second)
	}
	return DeadLinks// DeadLinks用于不重复地存放DeadLink
}



func main() {
	//GetResCode("https://www.javaroad.cn/questions/327057")
	//GetChildLink("https://www.javaroad.cn/questions/327057")
	//help("https://clslaid.icu/about")
	//help("https://www.javaroad.cn/questions/327057")
	DeadLinkDetector("https://clslaid.icu/about")
	//DeadLinkDetector("https://www.javaroad.cn/questions/327057")
	//DeadLinkDetector("https://www.ft.com/")
}