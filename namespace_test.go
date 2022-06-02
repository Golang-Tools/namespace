package namespace

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_namespace_to_str(t *testing.T) {
	namespace := NameSpcae{"a", "b", "c"}
	ns := namespace.ToString()
	assert.Equal(t, "a::b::c", ns)
}

func Test_namespace_to_str_with_delimiter(t *testing.T) {
	namespace := NameSpcae{"a", "b", "c"}
	ns := namespace.ToString(WithNamespaceDelimiter("??"))

	assert.Equal(t, "a??b??c", ns)
}

func Test_namespace_genkey(t *testing.T) {
	namespace := NameSpcae{"a", "b", "c"}
	k := namespace.FullName("q")
	assert.Equal(t, "a::b::c::q", k)
}

func Test_namespace_genkey_with_delimiter(t *testing.T) {
	namespace := NameSpcae{"a", "b", "c"}
	k := namespace.FullName("q", WithKeyDelimiter("-"))

	assert.Equal(t, "a::b::c-q", k)
}

func Test_namespace_fromfullname_same_delimiter(t *testing.T) {
	keyStr := "a::b::c"
	namespace, endpointStr, err := FromFullName(keyStr)
	if err != nil {
		assert.FailNow(t, err.Error(), "Gen namespace from key string get error")
	}
	assert.Equal(t, NameSpcae{"a", "b"}, namespace)
	assert.Equal(t, "c", endpointStr)
}
func Test_namespace_fromfullname_same_delimiter_no_namespace(t *testing.T) {
	keyStr := "c"
	_, k, err := FromFullName(keyStr)
	if err != nil {
		assert.Equal(t, ErrKeyNotHaveNamespace, err)
		assert.Equal(t, k, keyStr)
	}
}

func Test_namespace_fromfullname_with_delimiter_option(t *testing.T) {
	keyStr := "a:b::c"
	namespace, endpointStr, err := FromFullName(keyStr, WithNamespaceDelimiter(":"), WithKeyDelimiter("::"))
	if err != nil {
		assert.FailNow(t, err.Error(), "Gen namespace from key string get error")
	}
	assert.Equal(t, NameSpcae{"a", "b"}, namespace)
	assert.Equal(t, "c", endpointStr)
}

func Test_namespace_fromfullname_with_no_namespace(t *testing.T) {
	keyStr := "abc"
	_, k, err := FromFullName(keyStr, WithNamespaceDelimiter(":"), WithKeyDelimiter("::"))
	if err != nil {
		assert.Equal(t, ErrKeyNotHaveNamespace, err)
		assert.Equal(t, k, keyStr)

	} else {
		assert.FailNow(t, "can not get error")
	}
}

func Test_namespace_fromfullname_not_match(t *testing.T) {
	keyStr := "a::b:c:d"
	_, _, err := FromFullName(keyStr, WithKeyDelimiter(":"))
	if err != nil {
		assert.Equal(t, ErrNamespaceFormatNotMatch, err)
	} else {
		assert.FailNow(t, "can not get error")
	}
}

func Test_namespace_fromfullname_with_prixy_not_match(t *testing.T) {
	keyStr := "a::b::c:d"
	_, _, err := FromFullName(keyStr, WithKeyDelimiter(":"), WithPrefix("//"))
	if err != nil {
		assert.Equal(t, ErrNamespaceFormatNotMatch, err)
	} else {
		assert.FailNow(t, "can not get error")
	}
}

func Test_namespace_fromfullname_with_prixy(t *testing.T) {
	keyStr := "//a::b::c:d"
	_, _, err := FromFullName(keyStr, WithKeyDelimiter(":"), WithPrefix("//"))
	if err != nil {
		assert.Equal(t, ErrNamespaceFormatNotMatch, err)
	} else {
		assert.FailNow(t, "can not get error")
	}
}

func Test_namespace_to_str_with_prixy(t *testing.T) {
	namespace := NameSpcae{"a", "b", "c"}
	ns := namespace.ToString(WithPrefix("//"))

	assert.Equal(t, "//a::b::c", ns)
}

func Test_namespace_genkey_with_prixy(t *testing.T) {
	namespace := NameSpcae{"a", "b", "c"}
	k := namespace.FullName("q", WithKeyDelimiter("-"), WithPrefix("//"))

	assert.Equal(t, "//a::b::c-q", k)
}

func Test_namespace_genkey_with_redis_style(t *testing.T) {
	namespace := NameSpcae{"a", "b", "c"}
	k := namespace.FullName("q", WithRedisStyle())

	assert.Equal(t, "a::b::c::q", k)
}
func Test_namespace_fromfullname_with_redis_style(t *testing.T) {
	keyStr := "a::b::c::d"
	namespace, endpointStr, err := FromFullName(keyStr, WithRedisStyle())
	if err != nil {
		assert.FailNow(t, err.Error(), "Gen namespace from key string get error")
	}
	assert.Equal(t, NameSpcae{"a", "b", "c"}, namespace)
	assert.Equal(t, "d", endpointStr)
}

func Test_namespace_genkey_with_etcd_style(t *testing.T) {
	namespace := NameSpcae{"a", "b", "c"}
	k := namespace.FullName("q", WithEtcdStyle())

	assert.Equal(t, "/a/b/c/q", k)
}
func Test_namespace_fromfullname_with_etcd_style(t *testing.T) {
	keyStr := "/a/b/c/d"
	namespace, endpointStr, err := FromFullName(keyStr, WithEtcdStyle())
	if err != nil {
		assert.FailNow(t, err.Error(), "Gen namespace from key string get error")
	}
	assert.Equal(t, NameSpcae{"a", "b", "c"}, namespace)
	assert.Equal(t, "d", endpointStr)
}

func Test_SetDefaultOptions(t *testing.T) {
	SetDefaultOptions(WithEtcdStyle())
	keyStr := "/a/b/c/d"
	namespace, endpointStr, err := FromFullName(keyStr)
	if err != nil {
		assert.FailNow(t, err.Error(), "Gen namespace from key string get error")
	}
	assert.Equal(t, NameSpcae{"a", "b", "c"}, namespace)
	assert.Equal(t, "d", endpointStr)
	ReSetDefaultOptions()
	_, k, err := FromFullName(keyStr)
	if err != nil {
		assert.Equal(t, ErrKeyNotHaveNamespace, err)
		assert.Equal(t, k, keyStr)
	}
}

func Test_Random_key(t *testing.T) {
	SetDefaultOptions(WithEtcdStyle())
	namespace := NameSpcae{"a", "b", "c"}
	ns, err := namespace.RandomKey()
	if err != nil {
		assert.FailNow(t, err.Error(), "etcd style namespace RandomKey get error")
	}
	assert.Equal(t, true, strings.HasPrefix(ns, "/a/b/c/"))
	SetDefaultOptions(WithRedisStyle())
	ns, err = namespace.RandomKey()
	if err != nil {
		assert.FailNow(t, err.Error(), "redis style namespace RandomKey get error")
	}
	assert.Equal(t, true, strings.HasPrefix(ns, "a::b::c::"))
	ReSetDefaultOptions()
}
