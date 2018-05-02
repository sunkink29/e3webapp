
export interface Block {
	Subject, Description: string
	CurSize, MaxSize: number 
	RoomNumber: number
	BlockOpen: boolean
}

export interface Teacher {
	ID: string
	Email, Name: string
	Blocks : Block[]
	Current: boolean
}