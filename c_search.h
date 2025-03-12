      
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "stdint.h"

#define likely(x) __builtin_expect(!!(x), 1)
#define unlikely(x) __builtin_expect(!!(x), 0)

/*
    总占用内存为3G
*/

#define MAX_NODE_NUM 3000000
#define CHARSET_SIZE 255  /* 下一个字节最多只有255种可能 */

// 状态节点结构体
typedef struct {
    uint32_t next[CHARSET_SIZE];  // 字符集的下一个状态
    uint32_t isEnd;        // 是否是结束状态, 使用位运算，减少内存占用
} DFA_Node;


// 初始化DFA节点
void initDFA(DFA_Node *dfa, int *nodeCount) {
    /*
        对于每个节点, 子节点的个数设置为255个
    */
    for (int i = 0; i < CHARSET_SIZE; i++) {
        dfa[*nodeCount].next[i] = -1;
        dfa[*nodeCount].isEnd = 0;
    }
    (*nodeCount)++;
}

// 添加敏感词到DFA
void addWord(DFA_Node *dfa, int *nodeCount, const uint8_t *word) {
    // 小于2个字节的不算敏感词, 必须大于3个字节
    size_t word_len = strlen((char *)word);
    if (word_len < 2) {
	    return ;
    }
   
    int currentState = 0;
    // 遍历整个字符串
    for (int i = 0; i < word_len; i++) {
        uint8_t c = *(word + i);  /*直接使用字节作为位置使用*/
        if (dfa[currentState].next[c] == -1) {
            initDFA(dfa, nodeCount);  // 对应的字节位置不存在, 增加一个节点
            dfa[currentState].next[c] = *nodeCount - 1; // 存储节点所在的index
        }
	
        if (i == word_len -1) {
             dfa[currentState].isEnd = dfa[currentState].isEnd | ((uint32_t)1 << c);
	    }
	
    	currentState = dfa[currentState].next[c];  // 设置新的节点
    }
}

// 检测输入文本中的敏感词, 记录检测路径
int checkSensitiveWords(DFA_Node *dfa, const uint8_t *text , int *result, uint16_t max_search_count) {
    uint16_t search_count = 0;
    int flag = 0;

    for (int j = 0; j < strlen((char *)text); j++){
        int currentState = 0;

        if (unlikely(search_count >= max_search_count)){
            break;
        }

        for (int i = j; text[i]; i++) {
            uint8_t c = text[i];

            if (dfa[currentState].next[c] == -1) {
                break;
            }

            if ((dfa[currentState].isEnd & ((uint32_t)1 << c)) > 0) {
                *(result + search_count*2) = j;
                *(result + search_count*2 + 1) = i;
                search_count++;
                flag = 1;
            }

            if (dfa[currentState].next[c] != -1) {
                currentState = dfa[currentState].next[c];
            }

        }
    }

    return flag; // 没有发现敏感词
}

/*
    创建DFA树
*/
DFA_Node *createDFA_Node_And_Init(const uint8_t **sensitive_words, int length){
    /*
        申请内存, 并初始化所有节点, 占用2GB的内存, 后续可以考虑使用
    */
    DFA_Node *dfa = (DFA_Node *) malloc(sizeof(DFA_Node) * MAX_NODE_NUM);

    // 初始化根节点
    int nodeCount = 0;
    initDFA(dfa, &nodeCount);

    // 添加敏感词
    for (uint32_t i = 0; i < length; i++){
        addWord(dfa, &nodeCount, sensitive_words[i]);
    }

    // 超出节点数量返回一个空指针
    if (nodeCount >= MAX_NODE_NUM){
        free(dfa);
        printf("最大的节点数量为: %d, 支持数量: %u, 超出请扩容\n", nodeCount, MAX_NODE_NUM);
        return NULL;
    }

    return dfa;
}

// 释放内存, 一定要释放
void DFA_free(DFA_Node *dfa){
    free(dfa);
}

// int main() {
//     // 创建
//     uint8_t *test_words[3] = {"bad", "eval", "你好"};
//     DFA_Node *dfa = createDFA_Node_And_Init(test_words, 3);

    
//     // 检测输入文本
//     uint8_t *text = "This is a example, 你好.";
//     if (checkSensitiveWords(dfa, text)) {
//         printf("敏感词检测：发现敏感词！\n");
//     } else {
//         printf("敏感词检测：未发现敏感词。\n");
//     }

//     return 0;
// }

    