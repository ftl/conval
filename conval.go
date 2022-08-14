/*
The package conval helps to evaluate the log files from amateur radio contests.
*/
package conval

type OperatorMode string

const (
	SingleOperator OperatorMode = "single"
	MultiOperator  OperatorMode = "multi"
)

type PowerMode string

const (
	HighPower PowerMode = "high"
	LowPower  PowerMode = "low"
	QRPPower  PowerMode = "qrp"
)
