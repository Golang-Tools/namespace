# namespace

[![Go Reference](https://pkg.go.dev/badge/github.com/Golang-Tools/namespace.svg)](https://pkg.go.dev/github.com/Golang-Tools/namespace)
[![Go Version](https://img.shields.io/badge/go-1.22%2B-blue.svg)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/Golang-Tools/namespace.svg)](https://github.com/Golang-Tools/namespace/releases)
[![License](https://img.shields.io/github/license/Golang-Tools/namespace.svg)](./LICENSE)

命名空间工具:给 redis、etcd 这类自身不带命名空间概念的组件人为划分业务空间,负责**拼装**与**解析**带命名空间的键。

A tiny helper to build and parse namespaced keys (for redis, etcd and friends).

## 安装 / Install

```bash
go get github.com/Golang-Tools/namespace
```

需要 Go 1.22 及以上。

## 用法 / Usage

### 拼装命名空间

```golang
ns := namespace.Namespace{"a", "b", "c"}

ns.ToString()                           // a::b::c
ns.ToString(namespace.WithEtcdStyle())  // /a/b/c

ns.FullName("q")                        // a::b::c::q
ns.FullName("q", namespace.WithEtcdStyle())                                   // /a/b/c/q
ns.FullName("q", namespace.WithPrefix("//"), namespace.WithKeyDelimiter("-")) // //a::b::c-q
```

### 解析全名

```golang
ns, key, err := namespace.FromFullName("/a/b/c/d", namespace.WithEtcdStyle())
// ns == namespace.Namespace{"a", "b", "c"},key == "d"
```

`FullName` 与 `FromFullName` 使用同一套配置,可以闭环互转(详见 `TestFullNameRoundTrip`)。

### 生成随机键

```golang
key, err := ns.RandomKey(namespace.WithRedisStyle())               // a::b::c::<uuid4>
key, err = ns.RandomKey(namespace.WithRandomKeyGen(idgener.ULID))  // 换成 ulid 算法
```

### 调整默认配置

```golang
namespace.SetDefaultOptions(namespace.WithEtcdStyle())
defer namespace.ReSetDefaultOptions() //临时调整后记得复位

ns, key, err := namespace.FromFullName("/a/b/c/d")
```

- `SetDefaultOptions` 在已有默认值上叠加,未指定的字段保持不变;
- 传入 `With*` 选项的调用只影响本次调用,不会污染默认配置;
- 默认配置读取无锁、写入原子替换,并发使用是安全的,但仍建议在启动阶段完成配置。

## 配置项

| 选项 | 说明 |
| --- | --- |
| `WithPrefix(prefix)` | 命名空间前缀,如 etcd 风格的 `/` |
| `WithNamespaceDelimiter(delimiter)` | 各级命名空间之间的分隔符,默认 `::` |
| `WithKeyDelimiter(delimiter)` | 命名空间与 key 之间的分隔符,默认 `::` |
| `WithRandomKeyGen(algo)` | 随机键使用的算法(`idgener.UUID4`/`idgener.Sonyflake`/`idgener.Snowflake`/`idgener.ULID`),默认 `uuid4` |
| `WithRedisStyle()` | 预设:`::` 分隔、无前缀 |
| `WithEtcdStyle()` | 预设:`/` 分隔、前缀 `/` |

## 错误

| 错误 | 场景 |
| --- | --- |
| `ErrKeyNotHaveNamespace` | 全名只有一段(找不到分隔符),同时把原字符串作为 key 返回 |
| `ErrNamespaceFormatNotMatch` | 前缀不匹配、异名分隔符下多于两段、出现空段落或空 key |
| `ErrInvalidDelimiter` | 命名空间分隔符或 key 分隔符为空 |

## 迁移到 v1.0.0

- 类型 `NameSpcae` 正名为 `Namespace`;`NameSpcae` 保留为别名并标记 Deprecated,下游无需改动即可编译。
- 修复了选项静默失效与 `SetDefaultOptions` 无效的问题(源于 optparams v1.0.0 的拷贝语义),现在选项只作用于本次调用。
- `FromFullName` 更严格:空段落(`a::b::::c`)与空 key(`/a/b/c/`)会返回 `ErrNamespaceFormatNotMatch`。
- `FullName`/`RandomKey` 在命名空间或 key 为空时不再输出多余分隔符(此前得到 `::q`,现在得到 `q`)。
- 依赖升级:idgener v1.0.0、optparams v1.0.0、testify v1.12.1,go 版本底限 1.22。

## 开发

```bash
bash scripts/check.sh
```

gofmt / go vet / go test / go test -race / Example / benchmark smoke 一次跑完。本仓库不依赖 GitHub Actions。

## License

[MIT](./LICENSE)
