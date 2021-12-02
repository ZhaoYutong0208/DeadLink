package main

import (
	"container/list"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DeadLink 题目要求构造的死链结构，由死链url和上一级url组成
type DeadLink struct {
	DeadLink string
	SourceLink string
}

func DeadLinkDetector(rootUrl string) []DeadLink {
	// your implementation
	temp := make(map[string]DeadLink)
	DeadLinks := help(rootUrl, temp)
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

// help2 调用GetChildLink和visit
// 记得处理rootUrl不能访问的情况
func help(rootUrl string, DeadLinks map[string]DeadLink)  map[string]DeadLink{
	allLinks := make(map[string]interface{})
	allLinks[rootUrl] = struct{}{}
	queue := list.New()
	queue.PushBack(rootUrl)
	for queue.Len() != 0{
		for i:=0;i<=queue.Len()-1;i++{
			//我尝试在这边再次使用go func()并发，报内存溢出
			temp := queue.Front()
			processUrl := fmt.Sprintf("%v", temp.Value)
			//fmt.Printf("%s",processUrl)
			queue.Remove(temp)
			tempLinks1 := make(map[string]interface{})
			doc := GetChildLink(processUrl)
			processLinks := visit(tempLinks1, doc)// 得到processUrl下所有的链接列表，用[]string类型的links储存
			for link, _  := range processLinks{//待思考：link里有自己怎么办？
				go func() {
					resCode := GetResCode(link)
					//fmt.Printf("%s\r\n",link)
					//fmt.Printf("%d\r\n",resCode)
					if resCode < 200 || resCode >= 300 { // 情况1：判断为死链
						fmt.Printf("%s ", link)
						fmt.Printf("%d\r\n", resCode)
						// 借助temp，把死链加入DeadLinks
						var temp DeadLink// 待思考：是不是要delete它
						temp.DeadLink = link
						temp.SourceLink = processUrl
						DeadLinks[link] = temp
					}else{// 情况2：判断为活链
						_, ok := allLinks[link]// 情况2.1：allLinks里有link
						if !ok {
							allLinks[link] = struct{}{}
							queue.PushBack(link)
						}
					}
				}()
				time.Sleep(1 * time.Second)
			}
		}
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