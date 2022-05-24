package security

import (
	"github.com/gin-gonic/gin"
	"time"

	"bni.co.id/swift-filtering/database"
)

const SessionId = "X-Session-Id"

type Session struct {
	ID             string
	UserId         int
	RoleId         int
	OrganizationId int
	LogonDate      time.Time
	ExpireDate     time.Time
}

func NewAuthMiddleware() gin.HandlerFunc {
	var db = database.Open()

	return func(context *gin.Context) {
		sessionId := context.Request.Header.Get(SessionId)
		headers := context.Writer.Header()

		/*
			Access-Control-Allow-Origin: http://siteA.com
				Access-Control-Allow-Methods: GET, POST, PUT
				Access-Control-Allow-Headers: Content-Type
		*/
		headers.Add("Access-Control-Allow-Origin", "*")
		headers.Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		headers.Add("Access-Control-Allow-Headers", "Content-Type")
		headers.Add("Access-Control-Allow-Headers", SessionId)

		var session []Session
		db.Raw("SELECT id, user_id, role_id, organization_id, logon_date, expire_date FROM sys_session WHERE id = ?", sessionId).Scan(&session)

		// if already expired
		if len(session) == 0 || session[0].ExpireDate.Before(time.Now()) {
			context.AbortWithStatus(403)
		}
	}
}
