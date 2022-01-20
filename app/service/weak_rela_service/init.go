package weak_rela_service

import (
	"container/list"
	"edgelog/app/utils/observer_mode"
)

var SubjectHub1 *observer_mode.Subject

func init() {
	SubjectHub1 = &observer_mode.Subject{
		Observers: list.New(),
	}
	//  Start registering observer role business 
	obs1 := &observerSMS{}
	obs2 := &observerDeliver{}
	SubjectHub1.Attach(obs1)
	SubjectHub1.Attach(obs2)

}
