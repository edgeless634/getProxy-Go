# getProxy-Go
golang学习用，自动爬取国内代理

## 依赖
`htmlquery "github.com/antchfx/xquery/html"`

## 使用
参考main函数
```golang
func main() {
	channel := make(chan string, 10)
	go testedProxy(channel)
	for i := range channel {
		fmt.Println(i)
	}
}

```