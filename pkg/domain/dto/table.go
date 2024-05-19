package dto

type Column struct {
	Name    string
	Type    string
	IsList  bool
	PK      bool
	Comment string
}

type Columns []*Column
