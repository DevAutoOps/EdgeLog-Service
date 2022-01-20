package weak_rela_service

import (
	"edgelog/app/utils/observer_mode"
	"fmt"
)

// Simulate a business with weak relationship with the main business ， for example ： Send SMS 
type observerSMS struct {
}

func (c *observerSMS) Update(subject *observer_mode.Subject) {
	fmt.Printf(" Simulate sending SMS ， Received parameters ：%v\n", subject.GetParams())
}
