package cancelPolicy

// ICancelPolicyDirectRepository キャンセルポリシー関連の直仕入れrepositoryのインターフェース
type ICancelPolicyDirectRepository interface {
	Update(propertyID int64, cancelPolicy string) error
}
