//+build ignore


//核心算法来自：https://bbs.kanxue.com/thread-283252-1.htm
//SSE2需要硬件支持，但有极个别设备不支持，故做一个纯软运行……
//①通过掩码的方式使其支持 前、中、后通配符 及半字节；
//②通过SSE2指令集找到特征码字节序列中的第一个不为'??'的元素后，后续的字节只比较不是'??'的特征字节，优化比较字节数。

#include <iostream>
#include <windows.h>

int main() {
    //模拟内存中的字节数据
    BYTE MemByte[] = {
            0x33, 0xC9, 0x89, 0x0D, 0xB4, 0x67, 0x92, 0x77, 0x89, 0x0D, 0xB8, 0x67,
            0x92, 0x77, 0x88, 0x08, 0x38, 0x48, 0x02, 0x74, 0x05, 0xE8, 0x94, 0xFF,
            0xFF, 0xFF, 0x33, 0xC0, 0xC3, 0x8B, 0xFF, 0x55, 0x8B, 0xEC, 0x83, 0xE4, 0xF8,
    };

    //特征码为：?9 ?? 0? ?? 67
    //会处理成：F9 FF 0F FF 67 进行匹配
    std::string pattern = "?9 ?? 0? ?? 67";
    int index = 0;
    while ((index = pattern.find(' ', index)) >= 0) pattern.erase(index, 1); //去除特征码所有空格
    size_t len = pattern.length() / 2; //计算特征码长度
    size_t nFirstMatch = len;  // 跳过头部??，记录第一次匹配的位置半字符或非??，用于优化搜索
    BYTE *pMarkCode = new BYTE[len];  // 存储转换后的特征码字节
    BYTE *pWildcard = new BYTE[len];  // 存储特征字符串中??、?(??=FF、?=F、非?=0) 通配符

    //处理特征码字符串，转换成字节数组
    for (size_t i = 0; i < len; i++) {
        std::string tmpStr = pattern.substr(i * 2, 2);
        if ("??" == tmpStr) // 是"??"的特征字符
        {
            tmpStr = "FF";
            pWildcard[i] = 0xFF;
        } else  // 不是"??"的特征字符
        {
            if ('?' == tmpStr[0]) // 左半字节为'?'
            {
                tmpStr[0] = 'F';
                pWildcard[i] = (0xF << 4);
            } else if ('?' == tmpStr[1]) // 右半字节为'?'
            {
                tmpStr[1] = 'F';
                pWildcard[i] = 0xF;
            } else {
                pWildcard[i] = 0x0;
            }
            if (nFirstMatch == len) nFirstMatch = i;
        }

        pMarkCode[i] = strtoul(tmpStr.c_str(), nullptr, 16);
    }

    //搜索内存，匹配特征码算法
    for (size_t m = 0; m < sizeof(MemByte); ++m) {
        if (!((MemByte[m] | pWildcard[nFirstMatch]) ^ pMarkCode[nFirstMatch])) //匹配上第一个字节
        {
            size_t offset = m - nFirstMatch; //记录偏移量
            for (size_t n = nFirstMatch; n < len; ++n) //匹配后续字节
            {
                if (offset > sizeof(MemByte) - len) break; //超出内存范围
                if (pWildcard[n] != 0xFF)  //后续字节是"??"的通配符，跳过，这句代码可以优化搜索
                    if ((MemByte[offset + n] | pWildcard[n]) ^ pMarkCode[n]) break; //匹配失败
                if (n + 1 == len) //匹配成功
                {
                    printf("%Ix\n", MemByte[m - nFirstMatch]);
                }
            }
        }
    }
    system("pause");
    return 0;
}

