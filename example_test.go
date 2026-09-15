package namespace

import (
	"fmt"
	"strings"
)

func ExampleNamespace_ToString() {
	ns := Namespace{"a", "b", "c"}
	fmt.Println(ns.ToString(WithRedisStyle()))
	fmt.Println(ns.ToString(WithEtcdStyle()))
	// Output:
	// a::b::c
	// /a/b/c
}

func ExampleNamespace_FullName() {
	ns := Namespace{"a", "b", "c"}
	fmt.Println(ns.FullName("q", WithRedisStyle()))
	fmt.Println(ns.FullName("q", WithEtcdStyle()))
	// Output:
	// a::b::c::q
	// /a/b/c/q
}

func ExampleFromFullName() {
	ns, key, err := FromFullName("/a/b/c/d", WithEtcdStyle())
	fmt.Println(ns, key, err)

	_, key, err = FromFullName("q")
	fmt.Println(key, err)
	// Output:
	// [a b c] d <nil>
	// q key not have namespace
}

func ExampleNamespace_RandomKey() {
	ns := Namespace{"a", "b"}
	key, err := ns.RandomKey(WithRedisStyle())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(strings.HasPrefix(key, "a::b::"), len(strings.TrimPrefix(key, "a::b::")))
	// Output:
	// true 36
}

func ExampleSetDefaultOptions() {
	//调整默认配置后,不显式传选项的调用都会使用新配置
	SetDefaultOptions(WithEtcdStyle())
	defer ReSetDefaultOptions()

	ns, key, err := FromFullName("/a/b/c/d")
	fmt.Println(ns, key, err)
	// Output:
	// [a b c] d <nil>
}
