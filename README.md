# namespace

命名空间工具,简化拼凑命名流程

很多软件如redis,自己不带命名空间.我们必须人为的构造命名空间来分隔业务.本工具就是用来做这个的.

## 用法

+ 构造命名空间

     ```golang
    namespace := NameSpcae{"a", "b", "c"}
    n := namespace.ToString(WithRedisStyle())
    //a::b::c
    k := namespace.FullName("q", WithRedisStyle())
    //a::b::c::q
    ```

+ 解析命名空间

    ```golang
    SetDefaultOptions(WithEtcdStyle())
    keyStr := "/a/b/c/d"
    namespace, key, err := FromFullName(keyStr)
    ReSetDefaultOptions()
    ```
