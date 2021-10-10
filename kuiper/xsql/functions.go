package xsql

import (
	"github.com/cloustone/pandas/kuiper/plugins"
	"strings"
)

// ONLY use NewFunctionValuer function to initialize
type FunctionValuer struct {
	funcPlugins *funcPlugins
}

//Should only be called by stream to make sure a single instance for an operation
func NewFunctionValuer(p *funcPlugins) *FunctionValuer {
	fv := &FunctionValuer{
		funcPlugins: p,
	}
	return fv
}

func (*FunctionValuer) Value(_ string) (interface{}, bool) {
	return nil, false
}

func (*FunctionValuer) Meta(_ string) (interface{}, bool) {
	return nil, false
}

var aggFuncMap = map[string]string{"avg": "",
	"count": "",
	"max":   "", "min": "",
	"sum":         "",
	"collect":     "",
	"deduplicate": "",
}

var funcWithAsteriskSupportMap = map[string]string{
	"collect": "",
	"count":   "",
}

var mathFuncMap = map[string]string{"abs": "", "acos": "", "asin": "", "atan": "", "atan2": "",
	"bitand": "", "bitor": "", "bitxor": "", "bitnot": "",
	"ceil": "", "cos": "", "cosh": "",
	"exp": "",
	"ln":  "", "log": "",
	"mod":   "",
	"power": "",
	"rand":  "", "round": "",
	"sign": "", "sin": "", "sinh": "", "sqrt": "",
	"tan": "", "tanh": "",
}

var strFuncMap = map[string]string{"concat": "",
	"endswith":    "",
	"format_time": "",
	"indexof":     "",
	"length":      "", "lower": "", "lpad": "", "ltrim": "",
	"numbytes":       "",
	"regexp_matches": "", "regexp_replace": "", "regexp_substr": "", "rpad": "", "rtrim": "",
	"substring": "", "startswith": "", "split_value": "",
	"trim":  "",
	"upper": "",
}

var convFuncMap = map[string]string{"concat": "", "cast": "", "chr": "",
	"encode": "",
	"trunc":  "",
}

var hashFuncMap = map[string]string{"md5": "",
	"sha1": "", "sha256": "", "sha384": "", "sha512": "",
}

var jsonFuncMap = map[string]string{
	"json_path_query": "", "json_path_query_first": "", "json_path_exists": "",
}

var otherFuncMap = map[string]string{"isnull": "",
	"newuuid": "", "tstamp": "", "mqtt": "", "meta": "",
}

func (fv *FunctionValuer) Call(name string, args []interface{}) (interface{}, bool) {
	lowerName := strings.ToLower(name)
	if _, ok := mathFuncMap[lowerName]; ok {
		return mathCall(name, args)
	} else if _, ok := strFuncMap[lowerName]; ok {
		return strCall(lowerName, args)
	} else if _, ok := convFuncMap[lowerName]; ok {
		return convCall(lowerName, args)
	} else if _, ok := hashFuncMap[lowerName]; ok {
		return hashCall(lowerName, args)
	} else if _, ok := jsonFuncMap[lowerName]; ok {
		return jsonCall(lowerName, args)
	} else if _, ok := otherFuncMap[lowerName]; ok {
		return otherCall(lowerName, args)
	} else if _, ok := aggFuncMap[lowerName]; ok {
		return nil, false
	} else {
		nf, fctx, err := fv.funcPlugins.GetFuncFromPlugin(name)
		if err != nil {
			return err, false
		}
		if nf.IsAggregate() {
			return nil, false
		}
		logger := fctx.GetLogger()
		logger.Debugf("run func %s", name)
		result, ok := nf.Exec(args, fctx)
		logger.Debugf("run custom function %s, get result %v", name, result)
		return result, ok
	}
}

func IsAggStatement(node Node) bool {
	var r = false
	WalkFunc(node, func(n Node) {
		if f, ok := n.(*Call); ok {
			if ok := isAggFunc(f); ok {
				r = true
				return
			}
		} else if d, ok := n.(Dimensions); ok {
			ds := d.GetGroups()
			if ds != nil && len(ds) > 0 {
				r = true
				return
			}
		}
	})
	return r
}

func isAggFunc(f *Call) bool {
	fn := strings.ToLower(f.Name)
	if _, ok := aggFuncMap[fn]; ok {
		return true
	} else if _, ok := strFuncMap[fn]; ok {
		return false
	} else if _, ok := convFuncMap[fn]; ok {
		return false
	} else if _, ok := hashFuncMap[fn]; ok {
		return false
	} else if _, ok := otherFuncMap[fn]; ok {
		return false
	} else if _, ok := mathFuncMap[fn]; ok {
		return false
	} else {
		if nf, err := plugins.GetFunction(f.Name); err == nil {
			if nf.IsAggregate() {
				//Add cache
				aggFuncMap[fn] = ""
				return true
			}
		}
	}
	return false
}

func HasAggFuncs(node Node) bool {
	if node == nil {
		return false
	}
	var r = false
	WalkFunc(node, func(n Node) {
		if f, ok := n.(*Call); ok {
			if ok := isAggFunc(f); ok {
				r = true
				return
			}
		}
	})
	return r
}

func HasNoAggFuncs(node Node) bool {
	if node == nil {
		return false
	}
	var r = false
	WalkFunc(node, func(n Node) {
		if f, ok := n.(*Call); ok {
			if ok := isAggFunc(f); !ok {
				r = true
				return
			}
		}
	})
	return r
}
