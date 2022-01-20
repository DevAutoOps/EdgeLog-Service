package observer_mode

//  Observer role （Observer） Interface 
type ObserverInterface interface {
	//  Receive status update message 
	Update(*Subject)
}
