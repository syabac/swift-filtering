package scheduler

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"bni.co.id/swift-filtering/config"
	"bni.co.id/swift-filtering/database"
	"bni.co.id/swift-filtering/myjson"
)

// StartSendResponseJob start Send Response Job back to GLOBS, run every 10-seconds
func StartSendResponseJob() {
	db = database.Open()
	go func() {
		for true {
			time.Sleep(10 * time.Second)
			log.Println("executing StartSendResponseJob...")
			for _, res := range getPendingData() {
				if !checkSentBackStatus(res.ID) {
					continue
				}

				now := time.Now()
				updatePendingData(res.ID, 2, nil)

				if sendResponse(res) {
					// TODO: need more strore response data
					updatePendingData(res.ID, 1, &now)
				} else {
					updatePendingData(res.ID, 0, nil)
				}
			}
		}
	}()

	log.Println("SendResponseJob started")
}

type responseData struct {
	ID              uuid.UUID `gorm:"type:uniqueidentifier"`
	Header          string
	FollowedUpValue string
}

var db *gorm.DB

func getPendingData() []responseData {
	var result []responseData
	db.Raw(`SELECT TOP 100 id, header, followed_up_value
		FROM swift_transactions
		WHERE is_suspected = 1
		AND is_followed_up IS NOT NULL
		AND has_sent_back = 0
		ORDER BY id`).Scan(&result)

	return result
}

func checkSentBackStatus(id uuid.UUID) bool {
	var num int
	db.Raw(`SELECT 1
		FROM swift_transactions
		WHERE id = ?`, id).
		Row().
		Scan(&num)

	return num == 0
}

// TODO: this is incomplete logic, need to know API server target specs
func sendResponse(resData responseData) bool {
	url := config.GetSetting("swift.globs.url")
	data := map[string]interface{}{
		"header": resData.Header,
		"value":  resData.FollowedUpValue,
	}

	jsonData := myjson.Encode(data)
	response, err := http.Post(url, "application/json", bytes.NewBufferString(jsonData))

	if err != nil {
		log.Printf("Error occured when invoking GLOBS API: %v\n", err)
		return false
	}

	responseText, _ := ioutil.ReadAll(response.Body)
	dataHolder := make(map[string]interface{})

	myjson.Decode(string(responseText), &dataHolder)

	return true
}

func updatePendingData(id uuid.UUID, hasSentBack int, sentBackDate *time.Time) {
	db.Exec(`UPDATE swift_transactions
		SET has_sent_back = ?,
			sent_back_date = ?
		WHERE id = ?`, hasSentBack, sentBackDate, id)
}
