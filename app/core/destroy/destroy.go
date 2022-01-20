package destroy

import (
	"edgelog/app/core/event_manage"
	"edgelog/app/global/consts"
	"edgelog/app/global/variable"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	//   Used for monitoring system signals 
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM) //  Monitor for possible exit signals 
		received := <-c                                                                           // Value in receive signal pipeline 
		variable.ZapLog.Warn(consts.ProcessKilled, zap.String(" Signal value ", received.String()))
		(event_manage.CreateEventManageFactory()).FuzzyCall(variable.EventDestroyPrefix)
		close(c)
		os.Exit(1)
	}()

}
