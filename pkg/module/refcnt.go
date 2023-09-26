package module

import (
	"bytes"
	"fmt"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io/ioutil"
	"os"
	"strconv"
)

func Loaded(module string) (loaded bool, err error) {
	_, err = os.Stat(fmt.Sprintf("/sys/module/%s", module))
	if err == nil {
		loaded = true
	} else if os.IsNotExist(err) {
		loaded = false
		err = nil
	} else {
		awesome_error.CheckErr(err)
		return
	}
	return
}

// RefCount
// https://www.kernel.org/doc/html/latest/admin-guide/abi-stable.html#symbols-under-sys-module
// https://elixir.bootlin.com/linux/v4.18.13/source/kernel/module.c#L943
// each time use module, refcnt++
// void __module_get(struct module *module)
//
//	{
//		if (module) {
//			preempt_disable();
//			atomic_inc(&module->refcnt);
//			trace_module_get(module, _RET_IP_);
//			preempt_enable();
//		}
//	}
//
// if refcnt==0, module is loaded
// if refcnt==-1, means module is unloaded
// if this file not exists, unable to know whether this module is loaded
// for overlay, refcnt means the number of overlay mounted
func RefCount(module string) (refcnt int, err error) {
	content, err := ioutil.ReadFile(fmt.Sprintf("/sys/module/%s/refcnt", module))
	if err != nil {
		return
	}
	content = bytes.Trim(content, "\n")
	refcnt, err = strconv.Atoi(string(content))
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}
