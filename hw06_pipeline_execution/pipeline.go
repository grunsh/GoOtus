package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func inDone(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				for range in {
				}
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- val:
				case <-done:
					for range in {
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
		inChan = stage(inDone(inChan, done))
	}
	return inChan
}
