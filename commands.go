package splunk

type ControlCommand string

const (
	ControlCommandPause           ControlCommand = "pause"
	ControlCommandUnpause         ControlCommand = "unpause"
	ControlCommandFinalize        ControlCommand = "finalize"
	ControlCommandCancel          ControlCommand = "cancel"
	ControlCommandTouch           ControlCommand = "touch"
	ControlCommandSetTTL          ControlCommand = "setttl"
	ControlCommandSetPriority     ControlCommand = "setpriority"
	ControlCommandEnablePreview   ControlCommand = "enablepreview"
	ControlCommandDisablePreview  ControlCommand = "disablepreview"
	ControlCommandSetWorkLoadPool ControlCommand = "setworkloadpool"
)
