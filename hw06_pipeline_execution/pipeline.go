package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		ch := make(Bi)

		go func(ch Bi, o Out) {
			defer close(ch)

			for {
				select {
				case <-done:
					return
				case v, ok := <-o:
					if !ok {
						return
					}
					ch <- v
				}
			}
		}(ch, out)

		out = stage(ch)
	}

	return out
}
