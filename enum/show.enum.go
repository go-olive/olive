package enum

type ShowTaskStatusID = uint32

var ShowTaskStatus = struct {
	Absent      ShowTaskStatusID
	Monitoring  ShowTaskStatusID
	Downloading ShowTaskStatusID
}{
	Absent:      0,
	Monitoring:  1,
	Downloading: 2,
}
