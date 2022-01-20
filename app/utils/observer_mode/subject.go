package observer_mode

import "container/list"

//  Observer Management Center （subject）
type Subject struct {
	Observers *list.List
	params    interface{}
}

// Register observer role 
func (s *Subject) Attach(observe ObserverInterface) {
	s.Observers.PushBack(observe)
}

// Delete observer role 
func (s *Subject) Detach(observer ObserverInterface) {
	for ob := s.Observers.Front(); ob != nil; ob = ob.Next() {
		if ob.Value.(*ObserverInterface) == &observer {
			s.Observers.Remove(ob)
			break
		}
	}
}

// Notify all observers 
func (s *Subject) Notify() {
	var l_temp *list.List = list.New()
	for ob := s.Observers.Front(); ob != nil; ob = ob.Next() {
		l_temp.PushBack(ob.Value)
		ob.Value.(ObserverInterface).Update(s)
	}
	s.Observers = l_temp
}

func (s *Subject) BroadCast(args ...interface{}) {
	s.params = args
	s.Notify()
}

func (s *Subject) GetParams() interface{} {
	return s.params
}
