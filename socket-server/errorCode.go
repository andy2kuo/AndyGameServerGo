package socketserver

type ErrorCode int32

// 回傳錯誤代碼整數型態
func (code ErrorCode) Int() int32 {
	return int32(code)
}

// 錯誤代碼
// 新增一筆時請同樣在最下方的errList新增對應敘述
const (

	/*-----------------------*/
	/*通用*/
	/*-----------------------*/

	OK      = ErrorCode(0)  // 沒問題
	UNKNOWN = ErrorCode(-1) // 未知錯誤

	/*-----------------------*/
	/*系統*/
	/*-----------------------*/

	SYSTEM_STOP_BY_SERVER = ErrorCode(-100) // 從Server停止
	SYSTEM_STOP_BY_OS     = ErrorCode(-101) // 從OS系統停止

	/*-----------------------*/
	/*憑證*/
	/*-----------------------*/

	TOKEN_EXPIRE    = ErrorCode(-200) // 憑證過期
	TOKEN_INVALID   = ErrorCode(-201) // 憑證無效
	TOKEN_FAILED    = ErrorCode(-202) // 憑證處理過程中出現錯誤
	TOKEN_TOO_EARLY = ErrorCode(-203) // 憑證有效時間還未到

	/*-----------------------*/
	/*權限*/
	/*-----------------------*/

	ACCESS_DENIED    = ErrorCode(-300) // 操作被拒絕
	PERMISSION_EXIST = ErrorCode(-301) // 權限重複

	/*-----------------------*/
	/*關閉連線原因*/
	/*-----------------------*/

	DISCONNECT_BY_SERVER              = ErrorCode(-400) // Server停止
	DISCONNECT_BY_OS                  = ErrorCode(-401) // OS發送停止指令
	DISCONNECT_BY_CLIENT_TIME_OUT     = ErrorCode(-402) // 客戶端連線超時
	DISCONNECT_BY_CLIENT_STOP         = ErrorCode(-403) // 客戶端連接中斷
	DISCONNECT_BY_NO_RESPONSE         = ErrorCode(-404) // 客戶端無回應
	DISCONNECT_BY_CLIENT_ID_DUPLICATE = ErrorCode(-405) // 客戶端編號重複
)