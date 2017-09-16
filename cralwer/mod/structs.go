package mod

type Mod struct {
}

type RequestParams struct {
	ModUrl string
}

type Thread struct {
	Title    string
	UserName string
	Num      uint
	Url      string
}

type Result struct {
	Parent       string
	PreviousPage string
	NextPage     string
	Threads      []Thread
}
