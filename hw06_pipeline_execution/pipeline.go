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
				// Без этого тест done блокирется на зпись в канал.
				// С этим, проваливается по времени. Я в тупике. Прощу помощи. 5-ю ночь без сна.
				for range readChan {
				}
				return
			case val, ok = <-readChan:
				if !ok {
					return
				}
				select {
				case out <- val:
				case <-done:
					for range readChan {
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
