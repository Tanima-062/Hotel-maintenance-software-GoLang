package cancelPolicy

// ICancelPolicyTlRepository キャンセルポリシー関連のTLrepositoryのインターフェース
type ICancelPolicyTlRepository interface {
	// Update キャンセルポリシー更新
	UpsertCancelPolicyTl(propertyID int64, cancelPolicy string) error
}
