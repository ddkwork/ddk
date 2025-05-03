clone=./clone
git lfs install

cd $clone

git clone --recursive -b dev https://github.com/HyperDbg/HyperDbg HyperDbgDev ##include "SDK/Headers/Symbols.h" HyperDbg\hyperdbg\include\SDK\Headers\Datatypes.h
git clone --recursive https://github.com/Air14/HyperHide                      #搬运桌面快捷方法的xmake.lua
git clone --recursive https://github.com/393686984/vt-driver                  #
git clone --recursive https://github.com/3526779568/vt-debuger                #
git clone --recursive https://github.com/Gbps/gbhv                            #删除info文件，关闭测试签名，工程配置有问题
git clone --recursive https://github.com/Zero-Tang/NoirVisor                  #
git clone --recursive https://github.com/wbenny/injdrv                        #
git clone --recursive https://github.com/wbenny/DetoursNT                     #
git clone --recursive https://github.com/tandasat/DdiMon                      #win11 wdk不支持win7，警告视为错误，cpp和链接器

# static const uint8_t* skip_prefixes(const uint8_t* first, const uint8_t* last) noexcept
# {
#   return std::find_if(
#     first,
#     last,
#     [&](uint8_t byte) {
#       static constexpr uint8_t skip_table[] = { 0xF2, 0xF3, 0x2E, 0x36, 0x3E, 0x26, 0x64, 0x65, 0x2E, 0x3E, 0x66, 0x67 };
#       return std::find(std::begin(skip_table), std::end(skip_table), byte) == std::end(skip_table);
#     }
#   );//bug add this
# }
#static constexpr auto，删除constexpr，c编译，，警告视为错误，cpp和链接器，关闭测试签名，还有一个是注释掉noexcept之类的中间关键字
git clone --recursive https://github.com/wbenny/hvpp 
