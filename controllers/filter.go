package controllers

import (
	"net/http"
	"strings"

	"bni.co.id/swift-filtering/scenario"
	"bni.co.id/swift-filtering/scoring"
	"bni.co.id/swift-filtering/swift"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

// FilterController controller that handling
type FilterController struct {
	trxRepo      *swift.Repository
	scenarioRepo *scenario.Repository
}

// NewFilter create new FilterController object
// and prepare dependencies
func NewFilter() *FilterController {
	return &FilterController{
		trxRepo:      swift.NewRepository(),
		scenarioRepo: scenario.NewRepository(),
	}
}

type resultChannel struct {
	RuleID   uuid.UUID `gorm:"type:uniqueidentifier"`
	TagName  string
	TagValue string
	Data     []scoring.CalculatedSanctionData
}

// CheckFilter process swift filter
func (uc *FilterController) CheckFilter(c *gin.Context) {
	swiftText := c.PostForm("swift")
	isTest := c.PostForm("is_test") == "TRUE"

	if swiftText == "" || strings.Trim(swiftText, " ") == "" {
		return
	}

	swiftMessage := swift.Parse(swiftText)
	mt := swiftMessage.Block2.MessageType
	scenarioRepo := scenario.NewRepository()
	rules := scenarioRepo.GetRulesByMessageType(mt)

	var result []resultChannel
	var dataChannel = make(chan resultChannel)
	var channelCount = 0

	for _, rule := range rules {
		for _, ruleItem := range *rule.RuleItems {
			tagValues, isExist := swiftMessage.Tags[ruleItem.TagName]

			if !isExist {
				continue
			}

			engine := scoring.GetScoringEngineByName(ruleItem.ScoringMethod)

			if engine == nil {
				continue
			}

			for _, tagValue := range tagValues {
				channelCount++
				go func() {
					var tmpResult []scoring.CalculatedSanctionData
					tmpResult = engine.Calculate(tagValue, ruleItem.MinimumScore)
					dataChannel <- resultChannel{
						RuleID:   rule.ID,
						Data:     tmpResult,
						TagName:  ruleItem.TagName,
						TagValue: tagValue,
					}
				}()
			}
		}
	}

	status := "PASS"
	isSuspected := 0

	for i := 0; i < channelCount; i++ {
		dc := <-dataChannel

		if len(dc.Data) > 0 {
			result = append(result, dc)
			isSuspected = 1
			status = "HIT"
		}
	}

	if !isTest {
		go func() {
			swiftID := uc.trxRepo.Save(swiftMessage, isSuspected)

			for _, item := range result {
				for _, itemDetail := range item.Data {
					suspc := &scenario.Suspected{
						SwiftID:       swiftID,
						RuleID:        item.RuleID,
						TagName:       item.TagName,
						TagValue:      item.TagValue,
						Score:         itemDetail.Score,
						SanctionValue: itemDetail.ItemValue,
						MatchedWord:   itemDetail.MatchedWord,
						SanctionType:  itemDetail.SanctionType,
					}
					uc.scenarioRepo.SaveSuspectedResult(suspc)
				}
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"swift":  swiftText,
		"result": result,
	})
}
