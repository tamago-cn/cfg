# cfg
	基于 go-ini 的简易配置管理模块

## 基本理念
	将每个模块的配置作为一个节(`section`)，有配置修改需求的模块可在模块的`init`方法中注册一个节，即可方便的使用配置文件修改配置，注册时的参数值将作为参数默认值
	注册的同时，将模块的初始化方法与析构方法传入，以便在进程启动时对模块做相应的初始化操作，进程正常退出时做必要的析构操作


