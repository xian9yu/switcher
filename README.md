## switcher v3
switcher是一个多功能的端口转发工具，支持转发本地的端口，支持正则表达式转发（实现端口复用）。

程序占用CPU和内存资源极低。

###v3版本，参考自: 

1.https://github.com/jamescun/switcher

2.https://github.com/jackyspy/switcher

3.https://github.com/crabkun/Switcher

## 使用方法
配置好目录下的config.json后，直接运行就行
## 配置
打开程序目录下的config.json，你会看到类似下面的内容

### 主结构

    {
      "log_level": "debug",
      "rules": [
        规则配置
      ]
    }

### 规则配置
    {
      "name": "test",
      "listen": "0.0.0.0:1234",
      "enable_regexp": false,
      "targets": [
        目标配置
      ]
    }
### 目标配置
    {
      "regexp": "正则表达式",
      "address": "127.0.0.1:80"
    }
### 字段解释
#### 主结构 
1. log_level代表日志等级，有info/debug/error可以选
1. rules是规则配置数组，看下面
#### 规则配置
1. name是这个规则的名字，为了在日志中区分不同规则，建议取不同的名字
2. listen是这个规则监听的地址，0.0.0.0:1234代表监听所有网卡的1234端口
3. enable_regexp为是否开启正则表达式模式，后面有解释
4. first_packet_timeout为等待客户端第一个数据包的超时时间(**毫秒**)，仅开启正则表达式模式后有效，后面有解释
5. targets为目标配置数组，看下面

#### 目标配置
目标配置只有一种模式：**正则模式**。  

**上面**规则配置的**enable_regexp**为true或false决定了这个目标配置是普通模式还是正则模式。  

**正则模式**，即上面的**enable_regexp**为**true**，程序会根据客户端第一个数据包来匹配正则表达式，匹配成功就转发到指定的目标。
为了防止客户端长时间不发第一个数据包，故可以通过上面的规则配置的**first_packet_timeout**字段来配置超时时间（毫秒）  

目标配置有两个字段：  
1.regexp字段在正则模式才有用，代表正则表达式。  
2.address字段代表要转发的目标地址和端口，可以是本地的地址，也可以是远程地址

## 示例配置

        {
          "name": "正则模式示例",
          "listen": "          "listen": "127.0.0.1:1234",
          "enable_regexp": true,
          "first_packet_timeout": 5000,
          "targets": [
            {
              "regexp": "^\\x16\\x03\\x01\\x02",
              "address": "127.0.0.1:80"
            },
            {
              "regexp": "^SSH",
              "address": "127.0.0.1:22"
            }
          ]
        }


上面的配置为正则模式规则，监听本机的1234端口。   
只要有HTTPS浏览器连接，就会转发到127.0.0.1:80。只要有SSH客户端连接，就会转发到127.0.0.1:22。            

## 当前已实现协议的正则表达式
|协议|正则表达式|
| --- | ---|
|HTTPS(SSL)|^\x16\x03\x01\x02|
|SSH|^SSH|
|MYSQL|^\x20\x00|

**复制到JSON中记得注意特殊符号**

## 注意事项
本工具的正则模式的原理是根据客户端建立好连接后第一个数据包的特征进行判断是什么协议，然后再中转到目标地址。
这种方式已知有两个缺陷：

1. 不支持连接建立之后服务器主动握手的协议，例如VNC，FTP，MYSQL…(本项目中使未知协议延迟5s判断为MySQL协议)。
2. SSH无法连接请更换最新版putty或MobaXterm，因为SSH本来属于服务器主动握手的协议，但有些软件遵守有些软件不遵守，所以请选择客户端主动握手的软件。


##遇到了问题？
欢迎提issue或Pull Request

##License
3-Clause "Modified" BSD Licence.

[License](LICENSE)
