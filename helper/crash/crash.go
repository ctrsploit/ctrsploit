package crash

type Crash interface {
	Crash() (err error)
}

func MakeContainerCrash(cs ...Crash) (err error) {
	for _, c := range cs {
		err = c.Crash()
		if err != nil {
			return
		}
	}
	return
}
