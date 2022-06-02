package namespace

import (
	"strings"

	"github.com/Golang-Tools/idgener"
	"github.com/Golang-Tools/optparams"
)

type Options struct {
	Prefix             string                 //命名空间前缀
	NamespaceDelimiter string                 //各级命名空间间的分隔符
	KeyDelimiter       string                 //命名空间和key间的分隔符
	RandomKeyGen       idgener.IDGENAlgorithm //锁key的随机生成算法.默认uuidv4
}

var defaultOptions = Options{
	NamespaceDelimiter: "::",
	KeyDelimiter:       "::",
	RandomKeyGen:       idgener.IDGEN_UUIDV4,
}

//WithPrefix 设置命名空间间的前缀
func WithPrefix(prefix string) optparams.Option[Options] {
	return optparams.NewFuncOption(func(o *Options) {
		o.Prefix = prefix
	})
}

//WithNamespaceDelimiter 设置命名空间间的分割符
func WithNamespaceDelimiter(delimiter string) optparams.Option[Options] {
	return optparams.NewFuncOption(func(o *Options) {
		o.NamespaceDelimiter = delimiter
	})
}

//WithKeyDelimiter 设置键间的分割符
func WithKeyDelimiter(delimiter string) optparams.Option[Options] {
	return optparams.NewFuncOption(func(o *Options) {
		o.KeyDelimiter = delimiter
	})
}

//WithRedisStyle 设置redis风格的命名空间设置
func WithRedisStyle() optparams.Option[Options] {
	return optparams.NewFuncOption(func(o *Options) {
		o.KeyDelimiter = "::"
		o.NamespaceDelimiter = "::"
		o.Prefix = ""
	})
}

//WithEtcdStyle 设置etcd风格的命名空间设置
func WithEtcdStyle() optparams.Option[Options] {
	return optparams.NewFuncOption(func(o *Options) {
		o.KeyDelimiter = "/"
		o.NamespaceDelimiter = "/"
		o.Prefix = "/"
	})
}

//NameSpcae 带命名空间的键
type NameSpcae []string

//ToString namespace转换为字符串
//@params opts ...optparams.Option[Options] 设置项
func (n *NameSpcae) ToString(opts ...optparams.Option[Options]) string {
	opt := defaultOptions
	optparams.GetOption(&opt, opts...)
	namespaceStr := strings.Join(*n, opt.NamespaceDelimiter)
	if opt.Prefix == "" {
		return namespaceStr
	}
	builder := strings.Builder{}
	builder.WriteString(opt.Prefix)
	builder.WriteString(namespaceStr)
	return builder.String()
}

//FullName 在命名空间基础上创建一个key的全名
//@params key string 用于标识的键
//@params opts ...optparams.Option[Options] 设置项
func (n *NameSpcae) FullName(key string, opts ...optparams.Option[Options]) string {
	opt := defaultOptions
	optparams.GetOption(&opt, opts...)
	namespaceStr := strings.Join(*n, opt.NamespaceDelimiter)
	builder := strings.Builder{}
	if opt.Prefix != "" {
		builder.WriteString(opt.Prefix)
	}
	builder.WriteString(namespaceStr)
	builder.WriteString(opt.KeyDelimiter)
	builder.WriteString(key)
	return builder.String()
}

//RandomKey 在命名空间基础上创建一个随机的key的全名
//@params opts ...optparams.Option[Options] 设置项
func (n *NameSpcae) RandomKey(opts ...optparams.Option[Options]) (string, error) {
	opt := defaultOptions
	optparams.GetOption(&opt, opts...)
	namespaceStr := strings.Join(*n, opt.NamespaceDelimiter)
	randomkey, err := idgener.Next(opt.RandomKeyGen)
	if err != nil {
		return "", err
	}
	builder := strings.Builder{}
	if opt.Prefix != "" {
		builder.WriteString(opt.Prefix)
	}
	builder.WriteString(namespaceStr)
	builder.WriteString(opt.KeyDelimiter)
	builder.WriteString(randomkey)
	return builder.String(), nil
}

//FromFullName 从全名字符串中解析出命名空间和key
//@params fullname string 带解析全名
//@params opts ...optparams.Option[Options] 设置项
//@return NameSpcae 命名空间
//@return string key
//@return error 解析错误
func FromFullName(fullname string, opts ...optparams.Option[Options]) (NameSpcae, string, error) {
	opt := defaultOptions
	optparams.GetOption(&opt, opts...)
	if opt.Prefix != "" {
		if !strings.HasPrefix(fullname, opt.Prefix) {
			return nil, "", ErrNamespaceFormatNotMatch
		}
		fullname = fullname[len(opt.Prefix):]
	}
	nsinfo := strings.Split(fullname, opt.KeyDelimiter)
	len_nsinfo := len(nsinfo)
	if opt.KeyDelimiter == opt.NamespaceDelimiter {
		switch len_nsinfo {
		case 1:
			{
				return nil, nsinfo[0], ErrKeyNotHaveNamespace
			}
		default:
			{
				n := NameSpcae(nsinfo[:len_nsinfo-1])
				k := nsinfo[len_nsinfo-1]
				return n, k, nil
			}
		}
	} else {
		switch len_nsinfo {
		case 1:
			{
				return nil, nsinfo[0], ErrKeyNotHaveNamespace
			}
		case 2:
			{
				n := NameSpcae(strings.Split(nsinfo[0], opt.NamespaceDelimiter))
				k := nsinfo[1]
				return n, k, nil
			}
		default:
			{
				return nil, "", ErrNamespaceFormatNotMatch
			}
		}
	}
}

//SetDefaultOptions 设置默认命名空间配置
//@params opts ...optparams.Option[Options] 设置项
func SetDefaultOptions(opts ...optparams.Option[Options]) {
	optparams.GetOption(&defaultOptions, opts...)
}

//ReSetDefaultOptions 重置命名空间配置
func ReSetDefaultOptions() {
	defaultOptions = Options{
		NamespaceDelimiter: "::",
		KeyDelimiter:       "::",
	}
}

//WithRandomKeyGen 指定随机生成key时使用的随机算法
func WithRandomKeyGen(algo idgener.IDGENAlgorithm) optparams.Option[Options] {
	return optparams.NewFuncOption(func(o *Options) {
		o.RandomKeyGen = algo
	})
}
