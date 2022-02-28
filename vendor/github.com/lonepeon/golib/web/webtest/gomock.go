package webtest

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lonepeon/golib/web"
)

func MatchFlashErrorContains(v string) _GoMockFlashMessageContains {
	return _GoMockFlashMessageContains{kind: "error", value: v}
}

func MatchFlashSuccessContains(v string) _GoMockFlashMessageContains {
	return _GoMockFlashMessageContains{kind: "success", value: v}
}

func MatchDataContains(k string, v interface{}) _GoMockDataContains {
	return _GoMockDataContains{key: k, value: v}
}

type _GoMockFlashMessageContains struct {
	kind  string
	value string
}

func (m _GoMockFlashMessageContains) Matches(v interface{}) bool {
	flash, ok := v.(web.FlashMessage)
	if !ok {
		return false
	}

	return flash.Kind == m.kind && strings.Contains(flash.Message, m.value)
}

func (m _GoMockFlashMessageContains) String() string {
	return fmt.Sprintf("flash %s should contain %s", m.kind, m.value)
}

type _GoMockDataContains struct {
	key   string
	value interface{}
}

func (m _GoMockDataContains) Matches(v interface{}) bool {
	data, ok := v.(map[string]interface{})
	if !ok {
		return false
	}

	return reflect.DeepEqual(m.value, data[m.key])
}

func (m _GoMockDataContains) String() string {
	return fmt.Sprintf("key=%s, value=%#+v", m.key, m.value)
}
