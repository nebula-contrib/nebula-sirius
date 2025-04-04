package edge_alter

type AlterType string

const (
	AlterTypeAdd    AlterType = "ADD"
	AlterTypeDrop   AlterType = "DROP"
	AlterTypeChange AlterType = "CHANGE"
)

type IAlterEdgeTypeDefinition interface {
	GetAlterType() AlterType
	GenerateStatement() (string, error)
}
