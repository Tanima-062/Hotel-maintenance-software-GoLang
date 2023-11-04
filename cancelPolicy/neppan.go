package cancelPolicy

// ICancelPolicyNeppanRepository キャンセルポリシー関連のねっぱんrepositoryのインターフェース
type ICancelPolicyNeppanRepository interface {
	// Update キャンセルポリシー更新
	Update(propertyID int64, cancelPolicy string) error
}
