package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

var I int

func wrStage(in In, done In, s Stage) Out {
	out := make(Bi)
	var val interface{}
	var ok bool
	I++
	go func(stageNum int) {
		defer close(out)
		readChan := s(in)
		//fmt.Println("Запущен стейдж: ", stageNum)
		for {
			select {
			case <-done:
				//fmt.Println("Выход из стейдж по done: ", stageNum)
				for v := range readChan {
					_ = v
				}
				return
			case val, ok = <-readChan:
				if !ok {
					//fmt.Println("Выход из стейдж по закрытию канала чтения: ", stageNum)
					return
				}
				select {
				case out <- val:
				case <-done:
					//fmt.Println("Выход из стейдж по done2: ", stageNum)
					return
				}
			}
		}
	}(I)
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	inChan := in
	for _, stage := range stages {
		inChan = wrStage(inChan, done, stage)
	}
	return inChan
}
