package namespace

import (
	"errors"
)

//ErrKeyNotHaveNamespace key没有命名空间
var ErrKeyNotHaveNamespace = errors.New("key not have namespace")

//ErrNamespaceFormatNotMatch 命名空间格式不匹配
var ErrNamespaceFormatNotMatch = errors.New("namespace format not match")
