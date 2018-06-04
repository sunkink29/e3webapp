
export class Block {
	Subject: string
	Description: string
	CurSize: number
	MaxSize: number
    RoomNumber: number
    BlockOpen: boolean
}

export class Teacher {
	ID: string
	Email: string
	Name: string
	Blocks : Block[]
	Current: boolean
}

export class Student {
	ID: string
	Email: string 
	Name: string
	Grade: number
	Teacher1: string
	Teacher2: string
	Current: boolean
}