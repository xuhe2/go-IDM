package utils

type FilePart struct {
	Index        int
	From         int64
	To           int64
	Data         []byte
	FinishSignal chan int
}

func (fp *FilePart) Download() {

}
