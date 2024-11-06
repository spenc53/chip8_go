package chip

type CPU struct {
	PC         uint16
	I          uint16 // index register
	Stack      []uint16
	DelayTimer uint8
	SoundTimer uint8
	Reg        [16]uint8
}
