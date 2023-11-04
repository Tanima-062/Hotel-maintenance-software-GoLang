package cancelPolicy

// ICancelPolicyRaku2Repository キャンセルポリシー関連のらく通repositoryのインターフェース
type ICancelPolicyRaku2Repository interface {
	// Update キャンセルポリシー更新
	Update(propertyID int64, cancelPolicy string) error
}
