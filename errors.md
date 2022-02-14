# **关于error的一些处理**



## **一、目录**

+ Error vs Exception
+ Error Type
+ Handling Error
+ Go 1.13 errors
+ Go 2 Error Inspection

## 二、内容

### **1、Go中的error和其它语言的exception对比**

①、底层数据结构

```go
//标准库底下的error数据结构
type error interface {
    Error() string
}

//pkg/error下的基础数据结构
type fundamental struct {
	msg string//error信息
	*stack//堆栈指针
}

//Error接口实现
func (f *fundamental) Error() string { return f.msg }

```

②、error上抛返回指针而不是字符串

go中往往根据error携带的信息判断error种类，如果返回字符串可能导致由于重复命名error种类而导致没法正常判断

```go
type errorString string

func (e errorString) Error() string {
    return string(e)
}

func New(text string) error {
    return errorString(text)
}

var ErrNamedType = New("EOF")//自定义返回string
var ErrStructType = error.New("EOF")//标准库返回指针

func main() {
    if ErrNamedType == New("EOF") {
        fmt.Println("Named type Error")//字符串匹配导致判真
    }
    if ErrStructType == error.New("EOF") {
        fmt.Println("Struct Type Error")
    }
}
```

③、error的各语言演变史

+ C

单返回值，用指针做入参，返回值为int，表示成功或者失败

+ C++

引入了exception，但是无法知道被调用方会抛出什么异常

+ java

引入了checked exception，方法的所有者必须声明，调用者必须处理。在启动时大量抛出异常是正常行为。

+ go

go处理异常不是引用exception，它支持多参数返回，所以可以在函数签名中带上实现了error interface的对象，由调用者判断error类型。如果一个函数返回了(value, error)，必须先判定error。唯一可以忽略error的情况是，value无所谓。

go中的panic用于处理不可恢复的程序错误，例如索引越界、除零、不可恢复的环境问题、栈溢出。其他错误应用error判定。

