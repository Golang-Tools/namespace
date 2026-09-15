// Package namespace 用于构造和解析带命名空间的键。
//
// 命名空间是一组有序的字符串,配合分隔符拼装为键前缀或完整键名,
// 便于给 redis、etcd 这类自身不带命名空间概念的组件人为划分业务空间。
//
// 约定:
//   - 组合(ToString/FullName/RandomKey)与解析(FromFullName)使用同一套 Options,
//     因此 FullName 产出的全名可以被 FromFullName 解析回原始的命名空间和 key;
//   - 命名空间或 key 为空时不会输出多余的分隔符,避免产出悬挂分隔符的键名;
//   - 解析要求分隔符非空,且不接受空段落与空 key。
package namespace

import (
	"strings"
	"sync/atomic"

	"github.com/Golang-Tools/idgener"
	"github.com/Golang-Tools/optparams"
)

// 默认的分隔符
const (
	//DefaultNamespaceDelimiter 默认的命名空间层级分隔符
	DefaultNamespaceDelimiter = "::"
	//DefaultKeyDelimiter 默认的命名空间与key之间的分隔符
	DefaultKeyDelimiter = "::"
)

// Options 命名空间配置
type Options struct {
	Prefix             string            //命名空间前缀
	NamespaceDelimiter string            //各级命名空间间的分隔符
	KeyDelimiter       string            //命名空间和key间的分隔符
	RandomKeyGen       idgener.Algorithm //随机key的生成算法.默认uuidv4
}

// newDefaultOptions 构造一份完整的默认配置.
// 默认值与重置共用这一处定义,避免两边硬编码出现漂移.
func newDefaultOptions() Options {
	return Options{
		NamespaceDelimiter: DefaultNamespaceDelimiter,
		KeyDelimiter:       DefaultKeyDelimiter,
		RandomKeyGen:       idgener.UUID4,
	}
}

// defaultOptions 全局默认配置.
// 用 atomic.Pointer 承载:读取无锁、写入整体原子替换;配合 GetOption 的拷贝语义,
// 并发使用默认配置和调整默认配置之间不会产生数据竞争.
var defaultOptions atomic.Pointer[Options]

func init() {
	d := newDefaultOptions()
	defaultOptions.Store(&d)
}

// currentOptions 解析本次调用使用的配置.
// GetOption 是拷贝语义:基于全局默认配置拷贝后再叠加选项,既不会修改默认配置,
// 也不会被之前调用的选项污染,因此选项之间互不干扰.
// @params opts ...optparams.Option[Options] 设置项
func currentOptions(opts ...optparams.Option[Options]) *Options {
	return optparams.GetOption(defaultOptions.Load(), opts...)
}

// validate 校验配置是否可以用于拼接和解析
func (o Options) validate() error {
	if o.NamespaceDelimiter == "" || o.KeyDelimiter == "" {
		return ErrInvalidDelimiter
	}
	return nil
}

// WithPrefix 设置命名空间间的前缀
func WithPrefix(prefix string) optparams.Option[Options] {
	return optparams.NewFuncOption(func(o *Options) {
		o.Prefix = prefix
	})
}

// WithNamespaceDelimiter 设置命名空间间的分割符
func WithNamespaceDelimiter(delimiter string) optparams.Option[Options] {
	return optparams.NewFuncOption(func(o *Options) {
		o.NamespaceDelimiter = delimiter
	})
}

// WithKeyDelimiter 设置键间的分割符
func WithKeyDelimiter(delimiter string) optparams.Option[Options] {
	return optparams.NewFuncOption(func(o *Options) {
		o.KeyDelimiter = delimiter
	})
}

// WithRedisStyle 设置redis风格的命名空间设置
func WithRedisStyle() optparams.Option[Options] {
	return optparams.NewFuncOption(func(o *Options) {
		o.KeyDelimiter = DefaultKeyDelimiter
		o.NamespaceDelimiter = DefaultNamespaceDelimiter
		o.Prefix = ""
	})
}

// WithEtcdStyle 设置etcd风格的命名空间设置
func WithEtcdStyle() optparams.Option[Options] {
	return optparams.NewFuncOption(func(o *Options) {
		o.KeyDelimiter = "/"
		o.NamespaceDelimiter = "/"
		o.Prefix = "/"
	})
}

// Namespace 带命名空间的键
type Namespace []string

// NameSpcae 是 Namespace 的旧名(拼写有误).
//
// Deprecated: 使用 Namespace 替代.
type NameSpcae = Namespace

// namespaceString 把命名空间的各级拼接为字符串(不包含前缀)
func (n Namespace) namespaceString(opt *Options) string {
	return strings.Join(n, opt.NamespaceDelimiter)
}

// fullName 按配置把命名空间和key拼装为全名.
// 命名空间或key为空时不输出分隔符,避免产生无法回解析的悬挂分隔符.
func fullName(opt *Options, ns Namespace, key string) string {
	namespaceStr := ns.namespaceString(opt)
	builder := strings.Builder{}
	builder.Grow(len(opt.Prefix) + len(namespaceStr) + len(opt.KeyDelimiter) + len(key))
	builder.WriteString(opt.Prefix)
	builder.WriteString(namespaceStr)
	if namespaceStr != "" && key != "" {
		builder.WriteString(opt.KeyDelimiter)
	}
	builder.WriteString(key)
	return builder.String()
}

// ToString namespace转换为字符串
// @params opts ...optparams.Option[Options] 设置项
func (n Namespace) ToString(opts ...optparams.Option[Options]) string {
	opt := currentOptions(opts...)
	return opt.Prefix + n.namespaceString(opt)
}

// FullName 在命名空间基础上创建一个key的全名
// @params key string 用于标识的键
// @params opts ...optparams.Option[Options] 设置项
func (n Namespace) FullName(key string, opts ...optparams.Option[Options]) string {
	return fullName(currentOptions(opts...), n, key)
}

// RandomKey 在命名空间基础上创建一个随机的key的全名
// @params opts ...optparams.Option[Options] 设置项
// @return string 随机key的全名
// @return error 分隔符配置非法或随机算法不可用时返回错误
func (n Namespace) RandomKey(opts ...optparams.Option[Options]) (string, error) {
	opt := currentOptions(opts...)
	if err := opt.validate(); err != nil {
		return "", err
	}
	randomkey, err := idgener.Next(opt.RandomKeyGen)
	if err != nil {
		return "", err
	}
	return fullName(opt, n, randomkey), nil
}

// FromFullName 从全名字符串中解析出命名空间和key
// @params fullname string 待解析全名
// @params opts ...optparams.Option[Options] 设置项
// @return Namespace 命名空间
// @return string key,解析失败且能定位到key时返回去掉前缀后的原字符串
// @return error 解析错误:分隔符配置非法返回 ErrInvalidDelimiter;前缀不匹配、出现空段落或空key
// 返回 ErrNamespaceFormatNotMatch;全名中找不到分隔符返回 ErrKeyNotHaveNamespace
func FromFullName(fullname string, opts ...optparams.Option[Options]) (Namespace, string, error) {
	opt := currentOptions(opts...)
	if err := opt.validate(); err != nil {
		return nil, "", err
	}
	if opt.Prefix != "" {
		if !strings.HasPrefix(fullname, opt.Prefix) {
			return nil, "", ErrNamespaceFormatNotMatch
		}
		fullname = fullname[len(opt.Prefix):]
	}
	namespaceStr, key, err := splitFullName(fullname, opt)
	if err != nil {
		return nil, key, err
	}
	//解析结果不允许出现空段落与空key,保证 FromFullName(FullName(...)) 可以闭环
	ns := Namespace(strings.Split(namespaceStr, opt.NamespaceDelimiter))
	for _, segment := range ns {
		if segment == "" {
			return nil, "", ErrNamespaceFormatNotMatch
		}
	}
	if key == "" {
		return nil, "", ErrNamespaceFormatNotMatch
	}
	return ns, key, nil
}

// splitFullName 按配置把去掉前缀的全名拆分为命名空间字符串和key.
// 同名分隔符时以最后一个分隔符为界,之前的部分整体作为命名空间;异名分隔符时要求恰好两段.
// @return string 命名空间字符串(各级之间仍以命名空间分隔符相连,由调用方继续拆分)
// @return string key
// @return error 解析错误
func splitFullName(fullname string, opt *Options) (string, string, error) {
	if opt.KeyDelimiter == opt.NamespaceDelimiter {
		idx := strings.LastIndex(fullname, opt.KeyDelimiter)
		if idx < 0 {
			return "", fullname, ErrKeyNotHaveNamespace
		}
		return fullname[:idx], fullname[idx+len(opt.KeyDelimiter):], nil
	}
	nsinfo := strings.Split(fullname, opt.KeyDelimiter)
	switch len(nsinfo) {
	case 1:
		return "", nsinfo[0], ErrKeyNotHaveNamespace
	case 2:
		return nsinfo[0], nsinfo[1], nil
	default:
		return "", "", ErrNamespaceFormatNotMatch
	}
}

// SetDefaultOptions 设置默认命名空间配置
// 选项在已有默认配置的基础上叠加,未指定的字段保持不变;调整后所有不显式传选项的调用都会使用新配置.
// @params opts ...optparams.Option[Options] 设置项
func SetDefaultOptions(opts ...optparams.Option[Options]) {
	defaultOptions.Store(optparams.GetOption(defaultOptions.Load(), opts...))
}

// ReSetDefaultOptions 重置命名空间配置为内置默认值(分隔符"::"、无前缀、uuidv4随机算法)
func ReSetDefaultOptions() {
	d := newDefaultOptions()
	defaultOptions.Store(&d)
}

// WithRandomKeyGen 指定随机生成key时使用的随机算法
func WithRandomKeyGen(algo idgener.Algorithm) optparams.Option[Options] {
	return optparams.NewFuncOption(func(o *Options) {
		o.RandomKeyGen = algo
	})
}
