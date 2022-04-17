package enum

type StatusID = uint32

var Status = struct {
	Starting StatusID
	Pending  StatusID
	Running  StatusID
	Stopping StatusID
}{
	Starting: 0,
	Pending:  1,
	Running:  2,
	Stopping: 3,
}
