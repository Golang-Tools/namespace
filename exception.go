package namespace

import (
	"errors"
)

var (
	//ErrKeyNotHaveNamespace key没有命名空间
	ErrKeyNotHaveNamespace = errors.New("key not have namespace")

	//ErrNamespaceFormatNotMatch 命名空间格式不匹配.
	//包括前缀不匹配、异名分隔符下出现多于两段、空段落以及空key.
	ErrNamespaceFormatNotMatch = errors.New("namespace format not match")

	//ErrInvalidDelimiter 命名空间分隔符或key分隔符为空
	ErrInvalidDelimiter = errors.New("namespace delimiter or key delimiter is empty")
)
