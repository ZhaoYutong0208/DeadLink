package main
////使用BFS
//import (
//	"container/list"
//	"fmt"
//	"golang.org/x/net/html"
//	"io/ioutil"
//	"net/http"
//	"net/url"
//	"regexp"
//	"strings"
//	"sync"
//)
//
//// DeadLink 题目要求构造的死链结构，由死链url和上一级url组成
//type DeadLink struct {
//	DeadLink string
//	SourceLink string
//}
//
//// GetResCode 获取url网址的响应码resCode
//func GetResCode(rootUrl string) int{
//	u, _ := url.Parse(rootUrl)//将string解析成url格式
//	q := u.Query()//将路径解析为一个方便操作的对象
//	u.RawQuery = q.Encode()//将处理后url的回传给u
//	res, err := http.Get(u.String())
//	if err != nil {
//		fmt.Printf("get failed, err:%v\n", err)
//		return 0
//	}
//	resCode := res.StatusCode
//	return resCode
//}
//
//// GetChildLink 解析得到url网址的*html.Node，作为visit()访问遍历的输入
//func GetChildLink(rootUrl string) *html.Node{
//	u, _ := url.Parse(rootUrl)//将string解析成url格式
//	q := u.Query()//将路径解析为一个方便操作的对象
//	u.RawQuery = q.Encode()//将处理后url的回传给u
//	res, err := http.Get(u.String())
//	if err != nil {
//		fmt.Printf("get failed, err:%v\n", err)
//		return nil
//	}
//	body, err := ioutil.ReadAll(res.Body)//取出响应的body
//	if err != nil {
//		fmt.Printf("read from res.Body failed, err:%v\n", err)
//		return nil
//	}
//	b := string(body) // 转换成string类型
//	doc, err := html.Parse(strings.NewReader(b))
//	if err != nil {
//		fmt.Printf("findlinks: %v\n", err)
//		return nil
//	}
//	return doc
//}
//
//// visit 对*html.Node的内容迭代+遍历，得到该url下的网址列表
//// 参数说明：n *html.Node是要获得链接的对象
//// links map[string]interface{}无重复地返回该页面下所有满足条件的链接（活链+死链）
//// fixedUrl string正则表达式修正之后的根域名，用于筛选相同域名下的链接
//func visit(links map[string]interface{}, n *html.Node, fixedUrl string) map[string]interface{} {
//	// 用links这个map存储当前html下的网址
//	if n.Type == html.ElementNode && n.Data == "a" {
//		for _, a := range n.Attr {
//			if a.Key == "href" && strings.HasPrefix(a.Val, fixedUrl){
//				//resCode := GetResCode(a.Val)
//				links[a.Val] = struct{}{}
//			}
//		}
//	}
//	//这个递归，是对当前报文从上到下从里到外、逐级逐行检测
//	//FirstChild是第一个子节点，NextSibling是下一个同级节点
//	for c := n.FirstChild; c != nil; c = c.NextSibling {
//		links = visit(links, c, fixedUrl)
//	}
//	return links
//}
//
//
//// help 调用GetChildLink和visit，实现BFS的函数
//func help(rootUrl string, DeadLinks map[string]DeadLink)  map[string]DeadLink{
//	// Regular expression 正则表达式修正得到根域名fixedUrl
//	re := regexp.MustCompile("(\\w+):\\/\\/([^/:]+)")
//	fixedUrl := re.FindString(rootUrl)
//
//	allLinks := make(map[string]interface{})// allLinks用key存储所有处理过的活链
//	queue := list.New()// queue按层存储所有活链
//	queue.PushBack(rootUrl)
//	var wg sync.WaitGroup
//	var wg2 sync.WaitGroup
//	for queue.Len() != 0{
//		length := queue.Len()// 提前记录queue里存储的上一层活链数目
//		for i:=0;i<length;i++{
//			wg.Add(1) // 启动一个goroutine就登记+1
//			go func() {
//				defer wg.Done()
//				temp := queue.Front()
//				processUrl := fmt.Sprintf("%v", temp.Value)// processUrl 当前取出的待获取子链的网址
//				//fmt.Printf("%s",processUrl)
//				queue.Remove(temp)
//				tempMap := make(map[string]interface{})// 空map传入visit来存储返回的子链接
//				doc := GetChildLink(processUrl)// 获得url的html Node
//				processLinks := visit(tempMap, doc, fixedUrl)// 得到processUrl下所有的链接，用map类型的processLinks储存
//				//对processUrl下的链接进行判断
//				//判断1：是否在allLinks中，不是的话进行判断2，是的话不处理
//				//判断2：是否为死链，是的话加入到DeadLinks中，不是的话入队并加入allLinks中
//				for link, _  := range processLinks{
//					wg2.Add(1)
//					go func() {
//						defer wg2.Done()
//						_, ok := allLinks[link]
//						if !ok {//判断1：是否在allLinks中，不是的话进行判断2，是的话不处理
//							resCode := GetResCode(link)
//							fmt.Printf("%s\r\n",link)
//							fmt.Printf("%d\r\n",resCode)
//							if resCode < 200 || resCode >= 300 { // 判断2：是否为死链，是的话加入到DeadLinks中
//								fmt.Printf("%s ", link)
//								fmt.Printf("%d\r\n", resCode)
//								// 借助temp，把死链加入DeadLinks
//								var temp DeadLink
//								temp.DeadLink = link
//								temp.SourceLink = processUrl
//								DeadLinks[link] = temp
//							}else{// 不是死链的话入队并加入allLinks中
//								queue.PushBack(link)
//								allLinks[link] = struct{}{}
//							}
//						}
//					}()
//					wg2.Wait()
//				}
//			}()
//			wg.Wait() // 等待所有登记的goroutine都结束
//		}
//	}
//	return DeadLinks// DeadLinks用于不重复地存放DeadLink
//}
//
//// DeadLinkDetector 将help()函数返回的map[string]DeadLink格式转换成[]DeadLink并输出
//func DeadLinkDetector(rootUrl string) []DeadLink {
//	// your implementation
//	temp := make(map[string]DeadLink)
//	DeadLinks := help(rootUrl, temp)
//	//ans := make([]DeadLink, 0)
//	var ans []DeadLink
//	for _, link := range DeadLinks {
//		//temp := link.
//		fmt.Printf("DeadLink: %s\r\n",link.DeadLink)
//		fmt.Printf("SourceLink: %s\r\n",link.SourceLink)
//		ans = append(ans,link)
//	}
//	return ans
//}
//
//func main() {
//	//GetResCode("https://www.javaroad.cn/questions/327057")
//	//GetChildLink("https://www.javaroad.cn/questions/327057")
//	//help("https://clslaid.icu/about")
//	//help("https://www.javaroad.cn/questions/327057")
//	//DeadLinkDetector("https://clslaid.icu/about")
//	//DeadLinkDetector("https://www.javaroad.cn/questions/327057")
//	DeadLinkDetector("https://www.ft.com/")
//}