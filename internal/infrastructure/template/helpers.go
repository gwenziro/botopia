package template

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"
)

// JSONFunc adalah helper untuk JSON encoding di template
func JSONFunc(v interface{}) (template.JS, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return template.JS(b), nil
}

// SafeJSFunc menjadikan string aman untuk JavaScript
func SafeJSFunc(s string) template.JS {
	return template.JS(s)
}

// SafeCSSFunc menjadikan string aman untuk CSS
func SafeCSSFunc(s string) template.CSS {
	return template.CSS(s)
}

// DictFunc membuat map untuk template
func DictFunc(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("DictFunc requires an even number of arguments")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("DictFunc keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// FormatDateFunc formats a date as string
func FormatDateFunc(t time.Time, layout string) string {
	return t.Format(layout)
}

// TrimFunc trims whitespace from string
func TrimFunc(s string) string {
	return strings.TrimSpace(s)
}

// RegisterHelpers mendaftarkan semua helper function ke template engine
func RegisterHelpers(engine interface{}) {
	if htmlEngine, ok := engine.(*template.Template); ok {
		htmlEngine.Funcs(template.FuncMap{
			"json":       JSONFunc,
			"safeJS":     SafeJSFunc,
			"safeCSS":    SafeCSSFunc,
			"dict":       DictFunc,
			"formatDate": FormatDateFunc,
			"trim":       TrimFunc,
		})
	}
}
