# v1.0.0

首次正式发布:行为修复、依赖与命名的现代化。旧 API 以 Deprecated 别名保留,下游可直接编译。

## Fixed

- 修复选项静默失效:依赖 optparams v1.0.0 后 `GetOption` 变为拷贝语义,原有写法丢弃返回值,导致 `With*` 选项与 `SetDefaultOptions` 全部不生效。
- 修复 `ReSetDefaultOptions` 复位不完整:不再残留旧的 `RandomKeyGen`(此前仅因 uuid4 恰好是枚举零值而未暴露)。
- 修复空分隔符导致的错误解析:分隔符为空时 `strings.Split(s, "")` 会按 rune 拆分,现在 `FromFullName`/`RandomKey` 返回新增的 `ErrInvalidDelimiter`。
- 修复解析不拒绝空段落与空 key:`a::b::::c`、etcd 风格 `/a/b/c/` 此前会产出含空段的命名空间或空 key,现在返回 `ErrNamespaceFormatNotMatch`。
- 修复空命名空间与空 key 组合出悬挂分隔符:`Namespace{}.FullName("q")` 此前得到 `::q`,现在得到 `q`。
- 修复全局默认配置的并发读写竞争:改用 `atomic.Pointer[Options]`,读取无锁、写入整体原子替换。

## Changed

- 类型正名:`NameSpcae` → `Namespace`,`NameSpcae` 保留为别名并标记 Deprecated。
- 方法接收者由 `*NameSpcae` 改为值接收者,非寻址值也可以直接调用。
- `Options.RandomKeyGen` 与 `WithRandomKeyGen` 使用 `idgener.Algorithm`(即 `idgener.IDGENAlgorithm` 的现代化名称)。
- 新增 `DefaultNamespaceDelimiter`、`DefaultKeyDelimiter` 常量与 `ErrInvalidDelimiter` 哨兵错误。
- 依赖升级:idgener v1.0.0、optparams v1.0.0、testify v1.12.1;go 版本底限升至 1.22。

## Tooling

- 补充单元测试:选项隔离与默认值复位、边界与错误路径、往返解析、随机键唯一性、`-race` 并发用例,并新增 4 个 Benchmark。
- 新增 `example_test.go` 可执行示例。
- 新增本地校验脚本 `scripts/check.sh`(gofmt / go vet / go test / go test -race / Example / benchmark smoke),本仓库不依赖 GitHub Actions。
- 移除残留的 `pmfprc.json`;README 重写为中文为主的双语文档。

## Migration

- 直接升级即可:类型改名保留别名,`FullName`/`RandomKey`/`WithRedisStyle` 的输出格式不变,下游 redishelper 无需改动。
- `FromFullName` 更严格:空段落、空 key 会返回 `ErrNamespaceFormatNotMatch`,依赖旧的宽松行为时需要修正输入。
- `SetDefaultOptions` 语义不变(在已有默认值上叠加),但选项不再污染默认配置,带选项调用的结果更可预期。

# v0.0.2

## 新增接口

+ 新增方法`RandomKey`用于在命名空间基础上产生随机全名作为键

# v0.0.1

项目创建
