package cache

import "errors"

var ErrKeyNotFound = errors.New("key not found")
var ErrDecode = errors.New("decode error")
var ErrEncode = errors.New("encode error")
