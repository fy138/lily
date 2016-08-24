package gui

import (
	"encoding/json"
	"io/ioutil"
	"os"

	. "github.com/lxn/walk/declarative"
	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/conf"
	"github.com/tinycedar/lily/core"
	"github.com/tinycedar/lily/model"
)

func newTreeView() TreeView {
	treeModel := conf.Config.HostConfigModel
	return TreeView{
		AssignTo: &(context.treeView),
		Model:    treeModel,
		// click
		OnCurrentItemChanged: func() {
			context.deleteButton.SetEnabled(true)
			current := context.treeView.CurrentItem().(*model.HostConfigItem)
			if bytes, err := ioutil.ReadFile("conf/hosts/" + current.Text() + ".hosts"); err == nil {
				context.hostConfigText.SetText(string(bytes))
			} else {
				context.hostConfigText.SetText("")
			}
		},
		StretchFactor: 1,
		// double click
		OnItemActivated: func() {
			current := context.treeView.CurrentItem().(*model.HostConfigItem)
			previousIndex := conf.Config.CurrentHostIndex
			for i := 0; i < treeModel.RootCount(); i++ {
				item := treeModel.RootAt(i).(*model.HostConfigItem)
				if item.Text() == current.Text() {
					conf.Config.CurrentHostIndex = i
					item.Icon = common.IconMap[common.ICON_OPEN]
				} else if previousIndex == i {
					item.Icon = common.IconMap[common.ICON_NEW]
				}
				treeModel.PublishItemChanged(item)
			}
			if err := ioutil.WriteFile("C:/Windows/System32/drivers/etc/hosts", []byte(context.hostConfigText.Text()), os.ModeExclusive); err != nil {
				common.Error("Error writing to system hosts file: ", err)
			}
			configJSON, err := json.Marshal(conf.Config)
			if err != nil {
				common.Error("Error marshal json: %v", err)
			} else {
				ioutil.WriteFile("conf/config.json", configJSON, os.ModeExclusive)
			}
			// fire hosts switch
			core.FireHostsSwitch()
		},
	}
}
