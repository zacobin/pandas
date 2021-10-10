package plugins

import (
	"io/ioutil"
	"path"

	"github.com/cloustone/pandas/kuiper/util"
	ini "gopkg.in/ini.v1"
)

var g_uiMsg map[string]*ini.File

func getMsg(language, section, key string) string {
	language += ".ini"
	if conf, ok := g_uiMsg[language]; ok {
		s := conf.Section(section)
		if s != nil {
			return s.Key(key).String()
		}
	}
	return ""
}
func (m *Manager) readUiMsgDir() error {
	g_uiMsg = make(map[string]*ini.File)
	confDir, err := util.GetConfLoc()
	if nil != err {
		return err
	}

	dir := path.Join(confDir, "multilingual")
	infos, err := ioutil.ReadDir(dir)
	if nil != err {
		return err
	}

	for _, info := range infos {
		fName := info.Name()
		util.Log.Infof("uiMsg file : %s", fName)
		fPath := path.Join(dir, fName)
		if conf, err := ini.Load(fPath); nil != err {
			return err
		} else {
			g_uiMsg[fName] = conf
		}
	}
	return nil
}
