package service

import (
	"sync"
)

var (
	o               sync.Once
	DefaultRecord   RecordService
	DefaultField    FieldService
	DefaultModule   ModuleService
	DefaultPage     PageService
	DefaultWorkflow WorkflowService
	DefaultJob      JobService
)

func Init() {
	o.Do(func() {
		DefaultRecord = Record()
		DefaultField = Field()
		DefaultModule = Module()
		DefaultPage = Page()
		DefaultWorkflow = Workflow()
		DefaultJob = Job()
	})
}
