package main

/*
	#cgo CFLAGS: -I.
	#include "c_search.h"
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"
	"unsafe"
)

type CSearch struct {
	dfaNode *C.DFA_Node

	// Padding to achieve 64-byte alignment
	//_ [56]byte // 56 bytes
	rwMutex sync.RWMutex
}

// MakeCSearch 创建对象, 约占用3GB内存
func MakeCSearch(allSensitiveWords []string) (*CSearch, error) {
	cSearch := &CSearch{}

	// Create a C array to hold the pointers to the strings
	cStrings := make([]*C.uchar, len(allSensitiveWords))
	for i, s := range allSensitiveWords {
		cStrings[i] = (*C.uchar)(C.CBytes([]byte(s)))
		defer C.free(unsafe.Pointer(cStrings[i])) // 创建完成之后释放这块内存
	}

	// 调用 C 函数并得到一个指向 DFA_Node 的指针
	var dfaNodePtr *C.DFA_Node = C.createDFA_Node_And_Init(&cStrings[0], C.int(len(allSensitiveWords)))
	if dfaNodePtr == nil {
		return nil, fmt.Errorf("tree index build failed, please check and change MAX_NODE_NUM in c_search.h")
	}

	// 记录对象
	cSearch.dfaNode = dfaNodePtr

	return cSearch, nil
}

// UpdateCSearch 全量更新词表
func (s *CSearch) UpdateCSearch(allSensitiveWords []string) error {
	// 先拿锁, 避免多个协程同时更新造成内存浪费
	if ok := s.rwMutex.TryLock(); ok {
		defer s.rwMutex.Unlock()
	} else {
		return nil
	}

	// Create a C array to hold the pointers to the strings
	cStrings := make([]*C.uchar, len(allSensitiveWords))
	for i, s := range allSensitiveWords {
		cStrings[i] = (*C.uchar)(C.CBytes([]byte(s)))
		defer C.free(unsafe.Pointer(cStrings[i])) // 创建完成之后释放这块内存
	}

	// 调用 C 函数并得到一个指向 DFA_Node 的指针
	var dfaNodePtr *C.DFA_Node = C.createDFA_Node_And_Init(&cStrings[0], C.int(len(allSensitiveWords)))
	if dfaNodePtr == nil {
		return fmt.Errorf("tree index build failed, please check and change MAX_NODE_NUM in c_search.h")
	}

	// 保存旧的指针
	oldDFANodePtr := s.dfaNode

	// 替换掉旧的指针
	s.dfaNode = dfaNodePtr

	// 释放内存
	C.free(unsafe.Pointer(oldDFANodePtr))
	return nil
}

// CheckSensitive 查找一段句子中的敏感词
func (s *CSearch) CheckSensitive(text string, maxSearchCount uint16) []int {
	dfaNode := s.dfaNode

	// 拷贝这一段字符串内存到C
	searchString := (*C.uchar)(C.CBytes([]byte(text)))
	defer C.free(unsafe.Pointer(searchString))

	// 本地申请一段数组用于存储检索结果, 并将内容全部初始化为-1
	searchResult := make([]C.int, maxSearchCount)
	for i := range searchResult {
		searchResult[i] = -1
	}

	// 调用C进行检索
	C.checkSensitiveWords(dfaNode, searchString, (*C.int)(unsafe.Pointer(&searchResult[0])), C.ushort(maxSearchCount))

	// 将检索结果再转换会golang的格式
	var goSearchResult []int
	for _, index := range searchResult {
		if index != -1 {
			goSearchResult = append(goSearchResult, int(index))
		}
	}
	return goSearchResult
}

// DebugPrintSearchResult 展示提取到的敏感词结果, 用于调试
func (s *CSearch) DebugPrintSearchResult(text string, searchResult []int) {
	for i := 0; i < len(searchResult); i += 2 {
		start_index := searchResult[i]
		end_index := searchResult[i+1] + 1
		fmt.Printf("查找到的敏感词: %s, 开始和结束位置: %d -> %d \n", string([]byte(text)[start_index:end_index]), start_index, end_index)
	}

}

// example
func main() {
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

	// 全量更新敏感词
	err = cSearch.UpdateCSearch(allSensitiveWords)
	if err != nil {
		panic(err)
	}

	// 检索敏感词
	searchString := "习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平习近平"
	start := time.Now()                                                                // 获取当前时间
	searchResult := cSearch.CheckSensitive(searchString, 20)                           // 最多检索20个敏感词
	elapsed := time.Since(start)                                                       // 计算时间差
	fmt.Printf("总字符长度, %d, 函数调用时间: %d微妙\n", len(searchString), elapsed.Microseconds()) // 打印时间差

	// 打印查看结果
	cSearch.DebugPrintSearchResult(searchString, searchResult)
}
