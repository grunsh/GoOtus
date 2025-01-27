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
		defer func() {
			for range in {
			}
		}()

		for {
			select {
			case <-done:
				return
			default:
			}
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				out <- val
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	inChan := inDone(in, done)
	for _, stage := range stages {
		inChan = stage(inDone(inChan, done))
	}
	return inChan
}
