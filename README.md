# go-toy
🚀 Toy created by Golang

## 前言

刚开始学习go语言，看了基础的语法之后，就用一些小项目来练手，我认为，学习一门技术的最快，理解最深的就是敲了

## 组成

### 1 并行外部排序

通过这个课程可以初步感受到go routine之间的通信组件channel的使用

实现了：
1. 内部排序
2. 外部排序
3. 网络排序

充分利用了go routine的特点，对大数据集进行切分，分给不同的routine来进行处理，加速完成排序。

对于大数据集的排序，思想就是利用外部排序+归并

### 2 k-v存储

参考了rosedb作者的mini-db,在理解了之后，写出来这个toy,一方面是为了练习golang,一方面是为了理解rosedb,给作者点个赞。
同时我补充了部分单元测试和benchmark测试，使得可以更好的理解这个项目

mini-db采用了LSM的模型来构建一个可以持久化数据的k-v存储，在内存中建立了key对应的value在数据文件中的索引，也就是在文件中的偏移量


## 参考&感谢

- [imooc-搭建并行处理管道，感受GO语言魅力](https://www.imooc.com/learn/927)
- [roseduan-mindb](https://github.com/roseduan/minidb)