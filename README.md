### 微妙级别的敏感词提取

### 优势
* 1、静态内存，无动态内存分配
~~* 2、无锁化线程安全~~ 为避免不必要的麻烦，更新也是低频操作，这里还是加一把锁
* 3、微妙级别的速度

### 劣势
* 空间换时间，启动就会申请3GB内存

### 使用例子
* 拷贝c_search.go 和 c_search.h 到你的代码目录中即可
```go
    // 读取敏感词Json文件
	data, err := ioutil.ReadFile("sensitive_words.json")
	if err != nil {
		panic(err.Error())
	}
	var allSensitiveWords []string
	err = json.Unmarshal(data, &allSensitiveWords)
	if err != nil {
		panic(err.Error())
	}
    
	// 创建检索对象, 并初始化敏感词
	cSearch, err := MakeCSearch(allSensitiveWords)
	if err != nil {
		panic(err)
	}

	// 全量更新敏感词 thread safe
	err = cSearch.UpdateCSearch(allSensitiveWords)
	if err != nil {
		panic(err)
	}
	
	// 检索敏感词
	searchString := "敏1敏2敏感2test"
	start := time.Now()                                                                // 获取当前时间
	searchResult := cSearch.CheckSensitive(searchString, 20)                           // 最多检索20个敏感词
	elapsed := time.Since(start)                                                       // 计算时间差
	fmt.Printf("总字符长度, %d, 函数调用时间: %d微妙\n", len(searchString), elapsed.Microseconds()) // 打印时间差
    
	// 打印查看结果
	cSearch.DebugPrintSearchResult(searchString, searchResult)
```

### 联系作者
email发送至: liwenbo.csu@gmail.com