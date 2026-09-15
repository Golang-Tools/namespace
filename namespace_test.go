package namespace

import (
	"strings"
	"sync"
	"testing"

	"github.com/Golang-Tools/idgener"
	"github.com/Golang-Tools/optparams"
	"github.com/stretchr/testify/assert"
)

// 测试中反复使用的命名空间
var testNamespace = Namespace{"a", "b", "c"}

func TestNamespaceToString(t *testing.T) {
	tests := []struct {
		name      string
		namespace Namespace
		opts      []optparams.Option[Options]
		expected  string
	}{
		{name: "默认分隔符", namespace: testNamespace, expected: "a::b::c"},
		{name: "自定义命名空间分隔符", namespace: testNamespace, opts: []optparams.Option[Options]{WithNamespaceDelimiter("??")}, expected: "a??b??c"},
		{name: "带前缀", namespace: testNamespace, opts: []optparams.Option[Options]{WithPrefix("//")}, expected: "//a::b::c"},
		{name: "redis风格", namespace: testNamespace, opts: []optparams.Option[Options]{WithRedisStyle()}, expected: "a::b::c"},
		{name: "etcd风格", namespace: testNamespace, opts: []optparams.Option[Options]{WithEtcdStyle()}, expected: "/a/b/c"},
		{name: "空命名空间", namespace: Namespace{}, expected: ""},
		{name: "空命名空间带前缀", namespace: Namespace{}, opts: []optparams.Option[Options]{WithPrefix("//")}, expected: "//"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.namespace.ToString(tt.opts...))
		})
	}
}

func TestNamespaceFullName(t *testing.T) {
	tests := []struct {
		name      string
		namespace Namespace
		key       string
		opts      []optparams.Option[Options]
		expected  string
	}{
		{name: "默认分隔符", namespace: testNamespace, key: "q", expected: "a::b::c::q"},
		{name: "自定义key分隔符", namespace: testNamespace, key: "q", opts: []optparams.Option[Options]{WithKeyDelimiter("-")}, expected: "a::b::c-q"},
		{name: "前缀加自定义key分隔符", namespace: testNamespace, key: "q", opts: []optparams.Option[Options]{WithKeyDelimiter("-"), WithPrefix("//")}, expected: "//a::b::c-q"},
		{name: "redis风格", namespace: testNamespace, key: "q", opts: []optparams.Option[Options]{WithRedisStyle()}, expected: "a::b::c::q"},
		{name: "etcd风格", namespace: testNamespace, key: "q", opts: []optparams.Option[Options]{WithEtcdStyle()}, expected: "/a/b/c/q"},
		{name: "etcd风格自定义key分隔符", namespace: testNamespace, key: "q", opts: []optparams.Option[Options]{WithEtcdStyle(), WithKeyDelimiter(":")}, expected: "/a/b/c:q"},
		{name: "空命名空间", namespace: Namespace{}, key: "q", expected: "q"},
		{name: "空命名空间etcd风格", namespace: Namespace{}, key: "q", opts: []optparams.Option[Options]{WithEtcdStyle()}, expected: "/q"},
		{name: "空key", namespace: testNamespace, key: "", expected: "a::b::c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.namespace.FullName(tt.key, tt.opts...))
		})
	}
}

func TestFromFullName(t *testing.T) {
	tests := []struct {
		name        string
		fullname    string
		opts        []optparams.Option[Options]
		expectedNS  Namespace
		expectedKey string
		expectedErr error
	}{
		{name: "同名分隔符", fullname: "a::b::c", expectedNS: Namespace{"a", "b"}, expectedKey: "c"},
		{name: "异名分隔符", fullname: "a:b::c", opts: []optparams.Option[Options]{WithNamespaceDelimiter(":"), WithKeyDelimiter("::")}, expectedNS: Namespace{"a", "b"}, expectedKey: "c"},
		{name: "redis风格", fullname: "a::b::c::d", opts: []optparams.Option[Options]{WithRedisStyle()}, expectedNS: Namespace{"a", "b", "c"}, expectedKey: "d"},
		{name: "etcd风格", fullname: "/a/b/c/d", opts: []optparams.Option[Options]{WithEtcdStyle()}, expectedNS: Namespace{"a", "b", "c"}, expectedKey: "d"},
		{name: "带前缀", fullname: "//a::b::c::q", opts: []optparams.Option[Options]{WithPrefix("//")}, expectedNS: Namespace{"a", "b", "c"}, expectedKey: "q"},
		{name: "没有命名空间", fullname: "c", expectedKey: "c", expectedErr: ErrKeyNotHaveNamespace},
		{name: "异名分隔符下没有命名空间", fullname: "abc", opts: []optparams.Option[Options]{WithNamespaceDelimiter(":"), WithKeyDelimiter("::")}, expectedKey: "abc", expectedErr: ErrKeyNotHaveNamespace},
		{name: "异名分隔符下多于两段", fullname: "a::b:c:d", opts: []optparams.Option[Options]{WithKeyDelimiter(":")}, expectedErr: ErrNamespaceFormatNotMatch},
		{name: "前缀不匹配", fullname: "a::b::c:d", opts: []optparams.Option[Options]{WithKeyDelimiter(":"), WithPrefix("//")}, expectedErr: ErrNamespaceFormatNotMatch},
		{name: "同名分隔符下空段落", fullname: "a::b::::c", expectedErr: ErrNamespaceFormatNotMatch},
		{name: "etcd风格空key", fullname: "/a/b/c/", opts: []optparams.Option[Options]{WithEtcdStyle()}, expectedErr: ErrNamespaceFormatNotMatch},
		{name: "etcd风格空段落", fullname: "/a//b/c/d", opts: []optparams.Option[Options]{WithEtcdStyle()}, expectedErr: ErrNamespaceFormatNotMatch},
		{name: "空命名空间分隔符", fullname: "a::b::c", opts: []optparams.Option[Options]{WithNamespaceDelimiter("")}, expectedErr: ErrInvalidDelimiter},
		{name: "空key分隔符", fullname: "a::b::c", opts: []optparams.Option[Options]{WithKeyDelimiter("")}, expectedErr: ErrInvalidDelimiter},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns, key, err := FromFullName(tt.fullname, tt.opts...)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedKey, key)
			assert.Equal(t, tt.expectedNS, ns)
		})
	}
}

// TestFullNameRoundTrip FullName 产出的全名必须能被 FromFullName 还原
func TestFullNameRoundTrip(t *testing.T) {
	tests := []struct {
		name      string
		namespace Namespace
		key       string
		opts      []optparams.Option[Options]
	}{
		{name: "默认分隔符", namespace: testNamespace, key: "q"},
		{name: "redis风格", namespace: testNamespace, key: "q", opts: []optparams.Option[Options]{WithRedisStyle()}},
		{name: "etcd风格", namespace: testNamespace, key: "q", opts: []optparams.Option[Options]{WithEtcdStyle()}},
		{name: "前缀加异名key分隔符", namespace: testNamespace, key: "q", opts: []optparams.Option[Options]{WithPrefix("//"), WithKeyDelimiter("..")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullname := tt.namespace.FullName(tt.key, tt.opts...)
			ns, key, err := FromFullName(fullname, tt.opts...)
			assert.NoError(t, err)
			assert.Equal(t, tt.namespace, ns)
			assert.Equal(t, tt.key, key)
		})
	}
}

func TestNamespaceRandomKey(t *testing.T) {
	ns := testNamespace
	t.Run("etcd风格", func(t *testing.T) {
		key, err := ns.RandomKey(WithEtcdStyle())
		assert.NoError(t, err)
		assert.True(t, strings.HasPrefix(key, "/a/b/c/"))
	})
	t.Run("redis风格", func(t *testing.T) {
		key, err := ns.RandomKey(WithRedisStyle())
		assert.NoError(t, err)
		assert.True(t, strings.HasPrefix(key, "a::b::c::"))
	})
	t.Run("唯一性", func(t *testing.T) {
		seen := map[string]struct{}{}
		for i := 0; i < 100; i++ {
			key, err := ns.RandomKey()
			assert.NoError(t, err)
			_, ok := seen[key]
			assert.False(t, ok, "随机key出现重复: %s", key)
			seen[key] = struct{}{}
		}
	})
	t.Run("随机算法不可用", func(t *testing.T) {
		_, err := ns.RandomKey(WithRandomKeyGen(idgener.Algorithm(250)))
		assert.Error(t, err)
	})
	t.Run("空分隔符", func(t *testing.T) {
		_, err := ns.RandomKey(WithNamespaceDelimiter(""))
		assert.Equal(t, ErrInvalidDelimiter, err)
	})
}

// TestOptionsDoNotLeak 带选项的调用既不能失效,也不能污染默认配置
func TestOptionsDoNotLeak(t *testing.T) {
	ns := Namespace{"a", "b", "c"}
	assert.Equal(t, "/a/b/c/q", ns.FullName("q", WithEtcdStyle()))
	assert.Equal(t, "a::b::c::q", ns.FullName("q"))
	assert.Equal(t, "a-b-c::q", ns.FullName("q", WithNamespaceDelimiter("-")))
	assert.Equal(t, "a::b::c::q", ns.FullName("q"))
	assert.Equal(t, "//a::b::c", ns.ToString(WithPrefix("//")))
	assert.Equal(t, "a::b::c", ns.ToString())
	//选项按传入顺序叠加,后设置的生效
	assert.Equal(t, "/a/b/c:q", ns.FullName("q", WithEtcdStyle(), WithKeyDelimiter(":")))
	//解析同样不能污染默认配置
	_, _, err := FromFullName("/a/b/c/q", WithEtcdStyle())
	assert.NoError(t, err)
	_, _, err = FromFullName("a::b::c::q")
	assert.NoError(t, err)
}

func TestSetDefaultOptions(t *testing.T) {
	defer ReSetDefaultOptions()
	SetDefaultOptions(WithEtcdStyle())
	ns, key, err := FromFullName("/a/b/c/d")
	assert.NoError(t, err)
	assert.Equal(t, Namespace{"a", "b", "c"}, ns)
	assert.Equal(t, "d", key)

	ReSetDefaultOptions()
	_, key, err = FromFullName("/a/b/c/d")
	assert.Equal(t, ErrKeyNotHaveNamespace, err)
	assert.Equal(t, "/a/b/c/d", key)
}

// TestExplicitOptionsOverrideDefaults 显式传入的选项优先于全局默认配置
func TestExplicitOptionsOverrideDefaults(t *testing.T) {
	defer ReSetDefaultOptions()
	SetDefaultOptions(WithEtcdStyle())
	ns := Namespace{"a", "b", "c"}
	assert.Equal(t, "a::b::c::q", ns.FullName("q", WithRedisStyle()))
	assert.Equal(t, "/a/b/c/q", ns.FullName("q"))
}

// TestReSetDefaultOptionsRestoresRandomKeyGen 复位必须恢复随机算法(此前会残留旧算法)
func TestReSetDefaultOptionsRestoresRandomKeyGen(t *testing.T) {
	defer ReSetDefaultOptions()
	ns := Namespace{"a"}
	SetDefaultOptions(WithRandomKeyGen(idgener.ULID))
	key, err := ns.RandomKey()
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(key, "a::"))
	assert.Len(t, strings.TrimPrefix(key, "a::"), 26) //ULID为26字符

	ReSetDefaultOptions()
	key, err = ns.RandomKey()
	assert.NoError(t, err)
	assert.Len(t, strings.TrimPrefix(key, "a::"), 36) //uuid4为36字符
}

// TestNameSpcaeAlias 旧名别名保持可用,下游可以继续以 NameSpcae 构造
func TestNameSpcaeAlias(t *testing.T) {
	var ns NameSpcae = NameSpcae{"a", "b", "c"}
	var _ Namespace = ns
	assert.Equal(t, "a::b::c::q", ns.FullName("q", WithRedisStyle()))
	parsed, key, err := FromFullName("a::b::c::q", WithRedisStyle())
	assert.NoError(t, err)
	assert.Equal(t, Namespace(ns), parsed)
	assert.Equal(t, "q", key)
}

// TestValueReceiver 值接收者让非寻址的返回值也能直接调用
func TestValueReceiver(t *testing.T) {
	assert.Equal(t, "a::b::q", func() Namespace { return Namespace{"a", "b"} }().FullName("q"))
}

// TestConcurrentDefaultOptionsAndUse 并发调整默认配置与并发使用不应产生数据竞争(配合 go test -race)
func TestConcurrentDefaultOptionsAndUse(t *testing.T) {
	defer ReSetDefaultOptions()
	var wg sync.WaitGroup
	ns := Namespace{"a", "b", "c"}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				if i%2 == 0 {
					SetDefaultOptions(WithEtcdStyle())
				} else {
					SetDefaultOptions(WithPrefix("//"))
				}
				ReSetDefaultOptions()
			}
		}(i)
	}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = ns.FullName("q")
				_ = ns.ToString()
				_, _, _ = FromFullName("a::b::c::q")
			}
		}()
	}
	wg.Wait()
}

func BenchmarkNamespaceToString(b *testing.B) {
	ns := Namespace{"a", "b", "c", "d"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ns.ToString()
	}
}

func BenchmarkNamespaceFullName(b *testing.B) {
	ns := Namespace{"a", "b", "c", "d"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ns.FullName("q")
	}
}

func BenchmarkNamespaceFullNameWithOptions(b *testing.B) {
	ns := Namespace{"a", "b", "c", "d"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ns.FullName("q", WithRedisStyle())
	}
}

func BenchmarkFromFullName(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, _ = FromFullName("a::b::c::d::q", WithRedisStyle())
	}
}
