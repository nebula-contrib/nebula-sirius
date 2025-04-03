package tag_alter

type AlterType string

const (
	AlterTypeAdd    AlterType = "ADD"
	AlterTypeDrop   AlterType = "DROP"
	AlterTypeChange AlterType = "CHANGE"
)

type IAlterTypeDefinition interface {
	GetAlterType() AlterType
	GenerateStatement() (string, error)
}
