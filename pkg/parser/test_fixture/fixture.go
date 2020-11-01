package test_fixture

type TestModel struct {
	ID   []uint8 `json:"id" gots:"type:string"`
	Name string  `json:"name,omitempty"`
	Age  int     `json:"age" gots:"name:yearsAlive,optional"`
}
