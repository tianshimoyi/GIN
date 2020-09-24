package models

type Session struct {
	SessionId string `gorm:"primary_key;size:50"`     //sessionID用于标识唯一性
	UserName  string `gorm:"not null;unique;size:18"` //系统登录用户或管理员姓名
	UserId    uint   `gorm:"not null;unique"`         //系统登录用户或管理员ID
}

//查找session
func SearchSession(id string) Session {
	var session Session
	DB.First(&session, "session_id=?", id)
	return session
}

//插入session
func InsertSession(sid, name string, uid uint) error {
	var session Session
	DB.First(&session, "user_id=?", uid)
	if session.UserId > 0 {
		DB.Delete(&session)
	}
	return DB.Create(&Session{
		SessionId: sid,
		UserName:  name,
		UserId:    uid,
	}).Error
}

//删除session
func DeleteSession(session *Session) {
	DB.Delete(session)
}
