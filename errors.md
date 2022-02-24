# **关于error的一些处理**



## **一、目录**

+ Error vs Exception
+ Error Type
+ Handling Error
+ Go 1.13 errors
+ Go 2 Error Inspection

## 二、内容

### **1、error vs exception**

##### **①、底层数据结构**

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

##### **②、error上抛返回指针而不是字符串**

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

##### **③、error的各语言演变史**

+ C

单返回值，用指针做入参，返回值为int，表示成功或者失败

+ C++

引入了exception，但是无法知道被调用方会抛出什么异常

+ java

引入了checked exception，方法的所有者必须声明，调用者必须处理。在启动时大量抛出异常是正常行为。

+ go

go处理异常不是引用exception，它支持多参数返回，所以可以在函数签名中带上实现了error interface的对象，由调用者判断error类型。如果一个函数返回了(value, error)，必须先判定error。唯一可以忽略error的情况是，value无所谓。

go中的panic用于处理不可恢复的程序错误，例如索引越界、除零、不可恢复的环境问题、栈溢出。其他错误应用error判定。



### **2、error type**

##### **①、sentinel error**

特定错误，相当于把error定死一个格式或者某个值，以==判断不同error的不同处理方式。

当底层向上层返回更多的信息时，就会导致破坏本身的定义值，导致上层处理出错，很不方便。

如果公共函数或者方法返回一个sentinel error，则这个值必须是公有的，还需要有文档记录，一旦sentinel error多了，就会增加文档的表面积。

最糟糕的是这种模式会让包与包之间创建耦合，当我想检查错误是否是io.EOF时，就不得不引入io包（sentinel error在io包中定义）。

##### **②、errors types**

Error Type实现了error接口的自定义类型。

```go
type MyErr struct {
    Msg string
    File string
    Line int
}

//定义了错误内容，行号，文件等信息。
```

可以通过断言进行类型转换，调用者可以获取更多的上下文信息。

```go
func main() {
    err := test()
    switch err := err.(type) {
    case nil:
        //do something
    case *MyErr:
        fmt.Println("error occurred on line:", err.Line)
    default:
        //do something
    }
}
```

这种模式能让error携带更多信息，但由于需要断言判断，自定义的error类型就必须要成为公有变量，这会使得代码耦合度变高。

##### **③、opaque error**

透明化error

仅仅返回错误信息，处理交给上层业务。该模式较为推崇，但美中不足的是无法根据不同错误类型进行错误处理，也无法携带更多信息。



### **3、Handling Error**

应该只处理一次error，处理error意味着检查error类型，并作出单一的决定。

Go中的错误处理契约规定，在出错的情况下，不能对其他返回值的内容做出任何假设。由于JSON序列化失败，buf的内容是未知的，可能不包含任何内容，但也可能包含一个半写的JSON片段。

由于程序员在检查错误后记录但忘记返回，损坏的缓冲区被传递给上层，有可能配置文件被错误地写入，但函数返回的结果是正确的。

##### wrap errors注意事项

日志记录与错误无关的且对调试没有帮助的信息应该被视为噪声；

记录的原因是某些东西失败了，而日志包含了答案；

错误要被日志在顶层记录；

应用程序处理错误，保证100%完整性；

之后不再报告当前错误。



github.com/pkg/errors应用实例

```go
func ReadFile(path string) ([]byte, error) {
	f, errOpen := os.Open(path)
	if errOpen != nil {
		return nil, errors.Wrap(errOpen, "open failed")
	}
	defer f.Close()
	buf, errRead := ioutil.ReadAll(f)
	if errRead != nil {
		return nil, errors.Wrap(errRead, "read failed")
	}
	return buf, nil
}

func ReadConfig() ([]byte, error) {
	home := os.Getenv("HOME")
	config, err := ReadFile(filepath.Join(home, "setting.xml"))
	return config, errors.WithMessage(err, "could not read config")
}

func main() {
	_, err := ReadConfig()
	if err != nil {
        //errors.Cause返回根部错误，可用来sentienl error判断
		fmt.Printf("original err: %T %v\n", errors.Cause(err), errors.Cause(err))
        //%+v输出堆栈信息，先输出字段名字，再输出字段的值。%#v先输出结构体名字值，再输出内容
		fmt.Printf("stack trace:\n%+v\n", err)
		fmt.Println("err is:", err)
		os.Exit(1)
	}
}
```

在与其他库协作时，包括标准库，用errors.Wrap或者errors.Wrapf保存堆栈信息

直接返回错误，而不是在每个错误产生的地方都打出日志。在程序的顶部或者goroutine底部，使用%+v记录对战详情。

使用errors.Cause获取root error，再和sentinel error判定。

选择wrap error是只有applications可以选择应用的策略，具有最高可重用性的包只能返回根错误值。此机制与标准库使用的相同（kit库的sql.ErrNoRows）。

如果函数、方法不打算处理错误，那么用足够的上下文wrap errors并将其返回到调用堆栈中。

一旦确定函数、方法将处理错误，错误就不再是错误。如果仍需要返回，则不能返回错误值，只应返回零或者nil。

