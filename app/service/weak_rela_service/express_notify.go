package weak_rela_service

import (
	"edgelog/app/utils/observer_mode"
	"fmt"
)

// Simulate a business module that calls the logistics transporter interface to automatically create orders for the third party  ， It can be a separate file 
type observerDeliver struct {
	A int
}

func (c *observerDeliver) Update(subject *observer_mode.Subject) {
	fmt.Printf(" Simulate calling logistics transporter Api Interface ， Automatically notify the other party ：%v， %d\n", subject.GetParams(), c.A)
	c.A++
}
