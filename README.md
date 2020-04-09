基于LALR和链表的单词搜索
=======================================

使用
---------------------------------------
请按照：
```
$ go build -o main ./src
```
进行编译。因为本软件是基于Goolge开发的开源语言go编写而成，使用
前需要安装go的编译器。go语言是一门比较小众的语言，所以如果没有
该编译器的话，还请进入`https://golang.org/dl/`下载。

编译完成后，输入下命令可以得到主要用法：
```
$ ./main --help

Usage of ./main:
  -command string
    	the command to get the ID list (see README.pdf)
  -dirpath string
    	the input files' path (default "input")
  -mkindex
    	use this flag to make index named 'index.dict'
  -useindex
    	use file 'index.dict' to find result
```

通过`--mkindex`命令可以为`dirpath`文件夹（默认为input，但可以根
据`--dirpath`来进行修改）生成引索，然后通过`--useindex`来让其使
用引索文件。如果没有`--mkindex`和`--useindex`命令，则会动态遍历
文件夹，分词，以哈希表加链表作为数据结构（当然，使用引索文件的
话，虽然不需要遍历分词，但是还是需要抽象成数据结构）。

最后使用`--command`命令来完成主要目标：检索。语法见后文。

command实现与语法
---------------------------------------
这一次选择使用goyacc工具来生成相应的LALR解析器源码，我首先实现了
一个分词器，然后编写BNF式的语法规则，构建了解析器`parse.y`。

语法规则如下：
```
ast     ::= expr OR ast
         |  expr
expr    ::= atom AND expr
         |  atom
atom    ::= NOT STR
         |  STR
```

实例：
```
1. 'in' AND NOT 'in' OR 'her'
2. '\'' or not '\"'
3. 'outside' && !'inside'
```

如实例所述，STR字符串是由单引号或者双引号包裹而成，支持为单引号或
者双引号转义（其实在命令行结构下本身就需要对引号进行转义，所以难免
会发生二次转义，比如`--command="'\\\''"`这样比较丑陋的语法，但是
我的高数作业还有很多，实在是没有时间来支持交互式了）

而AND，OR，NOT不仅仅支持单词（并且忽略大小写）样式，还支持C式的逻辑
运算符，及`&&`，`||`，`!`。

（p.s.如果有时间的话可能会支持括号）

示例
---------------------------------------
``` shell
$ go build -o main ./src
$ ./main --help
Usage of ./main:
  -command string
    	the command to get the ID list (see README.pdf)
  -dirpath string
    	the input files' path (default "input")
  -mkindex
    	use this flag to make index named 'index.dict'
  -useindex
    	use file 'index.dict' to find result

$ ./main --command="'in' || 'not'"
result: [ d1.txt, d10.txt, d2.txt, d3.txt, d4.txt, d6.txt, d7.txt, d8.txt, d9.txt ]
$ ls
README.md
README.pdf
input
main
make.sh
src
$ ./main --mkindex
$ ls # it will make a new file named index.dict
README.md
README.pdf
index.dict
input
main
make.sh
src
$ ./main --useindex --command="not 'in'"
result: [ d5.txt ]
```
