package ceb


type CEBLogger struct {
	logChannel       chan LogLine
	isChannelCreated bool
}

type LogLine struct {
	Output string
	//MType  int
}

func (l CEBLogger) Log(Lable string, mType int) {

	line := LogLine{Output: Lable} //, MType: mType}

	if !l.isChannelCreated {
		l.isChannelCreated = true
		l.logChannel = make(chan LogLine)
		go consumeLogLines(l)
	}

	l.logChannel <- line

}

func consumeLogLines(l CEBLogger) {
	client := GetClient()

	select {
	case i := <-l.logChannel:
		client.ClientCommand(client.ListenerName, "log", "output", i)
	}
}
