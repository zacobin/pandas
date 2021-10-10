package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"strings"

	"github.com/cloustone/pandas/kuiper/util"
)

type (
	fileSource struct {
		About    *fileAbout              `json:"about"`
		Libs     []string                `json:"libs"`
		ConfKeys map[string][]*fileField `json:"properties"`
	}
	uiSource struct {
		About    *about              `json:"about"`
		Libs     []string            `json:"libs"`
		ConfKeys map[string][]*field `json:"properties"`
	}
	sourceProperty struct {
		cf   map[string]map[string]interface{}
		meta *uiSource
	}
)

func isInternalSource(fiName string) bool {
	internal := []string{`edgex.json`, `httppull.json`, `mqtt.json`}
	for _, v := range internal {
		if v == fiName {
			return true
		}
	}
	return false
}
func newUiSource(fi *fileSource) (*uiSource, error) {
	if nil == fi {
		return nil, nil
	}
	var err error
	ui := new(uiSource)
	ui.Libs = fi.Libs
	ui.About = newAbout(fi.About)
	ui.ConfKeys = make(map[string][]*field)

	for k, fields := range fi.ConfKeys {
		if ui.ConfKeys[k], err = newField(fields); nil != err {
			return nil, err
		}

	}
	return ui, nil
}

var g_sourceProperty map[string]*sourceProperty

func (m *Manager) uninstalSource(name string) {
	if v, ok := g_sourceProperty[name+".json"]; ok {
		if ui := v.meta; nil != ui {
			if nil != ui.About {
				ui.About.Installed = false
			}
		}
	}
}
func (m *Manager) readSourceMetaFile(filePath string) error {
	fileName := path.Base(filePath)
	if "mqtt_source.json" == fileName {
		fileName = "mqtt.json"
	}
	ptrMeta := new(fileSource)
	err := util.ReadJsonUnmarshal(filePath, ptrMeta)
	if nil != err || 0 == len(ptrMeta.ConfKeys) {
		return fmt.Errorf("file:%s err:%v", filePath, err)
	}
	if 0 == len(ptrMeta.ConfKeys["default"]) {
		return fmt.Errorf("not found default confKey %s", filePath)
	}
	if nil == ptrMeta.About {
		return fmt.Errorf("not found about of %s", filePath)
	} else if isInternalSource(fileName) {
		ptrMeta.About.Installed = true
	} else {
		_, ptrMeta.About.Installed = m.registry.Get(SOURCE, strings.TrimSuffix(fileName, `.json`))
	}

	yamlData := make(map[string]map[string]interface{})
	filePath = strings.TrimSuffix(filePath, `.json`) + `.yaml`
	err = util.ReadYamlUnmarshal(filePath, &yamlData)
	if nil != err {
		return fmt.Errorf("file:%s err:%v", filePath, err)
	}
	if 0 == len(yamlData["default"]) {
		return fmt.Errorf("not found default confKey from %s", filePath)
	}

	property := new(sourceProperty)
	property.cf = yamlData
	property.meta, err = newUiSource(ptrMeta)
	if nil != err {
		return err
	}

	g_sourceProperty[fileName] = property
	return err
}

func (m *Manager) readSourceMetaDir() error {
	g_sourceProperty = make(map[string]*sourceProperty)
	confDir, err := util.GetConfLoc()
	if nil != err {
		return err
	}

	dir := path.Join(confDir, "sources")
	infos, err := ioutil.ReadDir(dir)
	if nil != err {
		return err
	}

	if err = m.readSourceMetaFile(path.Join(confDir, "mqtt_source.json")); nil != err {
		return err
	}

	for _, info := range infos {
		fileName := info.Name()
		if strings.HasSuffix(fileName, ".json") {
			filePath := path.Join(dir, fileName)
			if err = m.readSourceMetaFile(filePath); nil != err {
				return err
			}
			util.Log.Infof("sourceMeta file : %s", fileName)
		}
	}
	return nil
}

func GetSourceConf(pluginName, language string) (b []byte, err error) {
	property, ok := g_sourceProperty[pluginName+".json"]
	if ok {
		cf := make(map[string]map[string]interface{})
		for key, kvs := range property.cf {
			aux := make(map[interface{}]interface{})
			for k, v := range kvs {
				aux[k] = v
			}
			cf[key] = util.ConvertMap(aux)
		}
		if b, err = json.Marshal(cf); nil != err {
			return nil, fmt.Errorf(`%s%v`, getMsg(language, source, "json_marshal_fail"), cf)
		} else {
			return b, err
		}
	}
	return nil, fmt.Errorf(`%s%s`, getMsg(language, source, "not_found_plugin"), pluginName)
}

func GetSourceMeta(pluginName, language string) (ptrSourceProperty *uiSource, err error) {
	property, ok := g_sourceProperty[pluginName+".json"]
	if ok {
		return property.cfToMeta(language)
	}
	return nil, fmt.Errorf(`%s%s`, getMsg(language, source, "not_found_plugin"), pluginName)
}

func GetSources() (sources []*pluginfo) {
	for fileName, v := range g_sourceProperty {
		node := new(pluginfo)
		node.Name = strings.TrimSuffix(fileName, `.json`)
		if nil == v.meta {
			continue
		}
		if nil == v.meta.About {
			continue
		}
		node.About = v.meta.About
		i := 0
		for ; i < len(sources); i++ {
			if node.Name <= sources[i].Name {
				sources = append(sources, node)
				copy(sources[i+1:], sources[i:])
				sources[i] = node
				break
			}
		}
		if len(sources) == i {
			sources = append(sources, node)
		}
	}
	return sources
}

func GetSourceConfKeys(pluginName string) (keys []string) {
	property := g_sourceProperty[pluginName+".json"]
	if nil == property {
		return keys
	}
	for k, _ := range property.cf {
		keys = append(keys, k)
	}
	return keys
}

func DelSourceConfKey(pluginName, confKey, language string) error {
	property := g_sourceProperty[pluginName+".json"]
	if nil == property {
		return fmt.Errorf(`%s%s`, getMsg(language, source, "not_found_plugin"), pluginName)
	}
	if nil == property.cf {
		return fmt.Errorf(`%s%s`, getMsg(language, source, "not_found_confkey"), confKey)
	}
	delete(property.cf, confKey)
	return property.saveCf(pluginName, language)
}

func AddSourceConfKey(pluginName, confKey, language string, content []byte) error {
	reqField := make(map[string]interface{})
	err := json.Unmarshal(content, &reqField)
	if nil != err {
		return err
	}

	property := g_sourceProperty[pluginName+".json"]
	if nil == property {
		return fmt.Errorf(`%s%s`, getMsg(language, source, "not_found_plugin"), pluginName)
	}

	if nil == property.cf {
		property.cf = make(map[string]map[string]interface{})
	}

	if 0 != len(property.cf[confKey]) {
		return fmt.Errorf(`%s%s`, getMsg(language, source, "confkey_already_exist"), pluginName)
	}

	property.cf[confKey] = reqField
	g_sourceProperty[pluginName+".json"] = property
	return property.saveCf(pluginName, language)
}

func AddSourceConfKeyField(pluginName, confKey, language string, content []byte) error {
	reqField := make(map[string]interface{})
	err := json.Unmarshal(content, &reqField)
	if nil != err {
		return err
	}

	property := g_sourceProperty[pluginName+".json"]
	if nil == property {
		return fmt.Errorf(`%s%s`, getMsg(language, source, "not_found_plugin"), pluginName)
	}

	if nil == property.cf {
		return fmt.Errorf(`%s%s`, getMsg(language, source, "not_found_confkey"), confKey)
	}

	if nil == property.cf[confKey] {
		return fmt.Errorf(`%s%s`, getMsg(language, source, "not_found_confkey"), confKey)
	}

	for k, v := range reqField {
		property.cf[confKey][k] = v
	}
	return property.saveCf(pluginName, language)
}

func recursionDelMap(cf, fields map[string]interface{}, language string) error {
	for k, v := range fields {
		if nil == v {
			delete(cf, k)
			continue
		}

		if delKey, ok := v.(string); ok {
			if 0 == len(delKey) {
				delete(cf, k)
				continue
			}

			var auxCf map[string]interface{}
			if err := util.MapToStruct(cf[k], &auxCf); nil != err {
				return fmt.Errorf(`%s%s.%s`, getMsg(language, source, "type_conversion_fail"), k, delKey)
			}
			cf[k] = auxCf
			delete(auxCf, delKey)
			continue
		}
		if reflect.Map == reflect.TypeOf(v).Kind() {
			var auxCf, auxFields map[string]interface{}
			if err := util.MapToStruct(cf[k], &auxCf); nil != err {
				return fmt.Errorf(`%s%s.%v`, getMsg(language, source, "type_conversion_fail"), k, v)
			}
			cf[k] = auxCf
			if err := util.MapToStruct(v, &auxFields); nil != err {
				return fmt.Errorf(`%s%s.%v`, getMsg(language, source, "type_conversion_fail"), k, v)
			}
			if err := recursionDelMap(auxCf, auxFields, language); nil != err {
				return err
			}
		}
	}
	return nil
}

func DelSourceConfKeyField(pluginName, confKey, language string, content []byte) error {
	reqField := make(map[string]interface{})
	err := json.Unmarshal(content, &reqField)
	if nil != err {
		return err
	}

	property := g_sourceProperty[pluginName+".json"]
	if nil == property {
		return fmt.Errorf(`%s%s`, getMsg(language, source, "not_found_plugin"), pluginName)
	}

	if nil == property.cf {
		return fmt.Errorf(`%s%s`, getMsg(language, source, "not_found_confkey"), confKey)
	}

	if nil == property.cf[confKey] {
		return fmt.Errorf(`%s%s`, getMsg(language, source, "not_found_confkey"), confKey)
	}

	err = recursionDelMap(property.cf[confKey], reqField, language)
	if nil != err {
		return err
	}
	return property.saveCf(pluginName, language)
}

func recursionNewFields(template []*field, conf map[string]interface{}, ret *[]*field) error {
	for i := 0; i < len(template); i++ {
		p := new(field)
		*p = *template[i]
		*ret = append(*ret, p)
		v, ok := conf[template[i].Name]
		if ok {
			p.Exist = true
		} else {
			p.Exist = false
			continue
		}

		var auxRet, auxTemplate []*field
		p.Default = &auxRet
		if nil == v {
			p.Default = v
		} else {
			if reflect.Map == reflect.TypeOf(v).Kind() {
				var nextCf map[string]interface{}
				if tmp, ok := v.(map[interface{}]interface{}); ok {
					nextCf = util.ConvertMap(tmp)
				} else {
					if err := util.MapToStruct(v, &nextCf); nil != err {
						return err
					}
				}
				if err := util.MapToStruct(template[i].Default, &auxTemplate); nil != err {
					return err
				}
				if err := recursionNewFields(auxTemplate, nextCf, &auxRet); nil != err {
					return err
				}
			} else {
				p.Default = v
			}
		}
	}
	return nil
}

func (this *sourceProperty) cfToMeta(language string) (*uiSource, error) {
	fields := this.meta.ConfKeys["default"]
	ret := make(map[string][]*field)
	for k, kvs := range this.cf {
		var sli []*field
		err := recursionNewFields(fields, kvs, &sli)
		if nil != err {
			return nil, fmt.Errorf(`%s%v`, getMsg(language, "source", "type_conversion_fail"), err)
		}
		ret[k] = sli
	}
	meta := new(uiSource)
	*meta = *(this.meta)
	meta.ConfKeys = ret
	return meta, nil
}

func (this *sourceProperty) saveCf(pluginName, language string) error {
	confDir, err := util.GetConfLoc()
	if nil != err {
		return fmt.Errorf(`%s%v`, getMsg(language, source, "not_found_file"), err)
	}

	dir := path.Join(confDir, "sources")
	if "mqtt" == pluginName {
		pluginName = "mqtt_source"
		dir = confDir
	}
	filePath := path.Join(dir, pluginName+".yaml")
	err = util.WriteYamlMarshal(filePath, this.cf)
	if nil != err {
		return fmt.Errorf(`%s%v`, getMsg(language, "source", "write_data_fail"), err)
	}
	return nil
}
