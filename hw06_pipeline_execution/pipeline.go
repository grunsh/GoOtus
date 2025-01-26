package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func wrStage(in In, done In, s Stage) Out {
	out := make(Bi)
	var val interface{}
	var ok bool
	go func() {
		defer close(out)
		readChan := s(in)
		for {
			select {
			case <-done:
				for v := range readChan {
					_ = v
				}
				return
			case val, ok = <-readChan:
				if !ok {
					return
				}
				select {
				case out <- val:
				case <-done:
					for v := range readChan {
						_ = v
					}
					return
				}
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	inChan := in
	for _, stage := range stages {
		inChan = wrStage(inChan, done, stage)
	}
	return inChan
}
