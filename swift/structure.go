package swift

import (
	"strings"
)

//Message swift message structures
type Message struct {
	Block1 *Block1
	Block2 *Block2
	Tags   map[string][]string
	Body   string
}

//Block1 Swift Block 1
type Block1 struct {
	Values          string
	ApplicationID   string
	ServiceID       string
	LogicalTerminal string
	SessionNumber   string
	SequenceNumber  string
}

// Block2Input contains information incoming transaction
type Block2Input struct {
	DestinationAddress string
	DeliveryMonitoring string
	ObsolescencePeriod string
}

// Block2Output contains information outcoming transaction
type Block2Output struct {
	InputTime          string
	MIRDate            string
	MIRLogicalTerminal string
	MIRSessionNumber   string
	MIRSequenceNumber  string
	OutputDate         string
	OutputTime         string
}

// Block2 Swift Block 2
type Block2 struct {
	Values            string
	IO                string
	MessageType       string
	MessagePriority   string
	InputDescription  *Block2Input
	OutputDescription *Block2Output
}

func (sb *Block1) loadValues(values string) {
	values = strings.ReplaceAll(values, " ", "")
	if values == "" {
		return
	}

	sb.Values = values
	sb.ApplicationID = string(values[0])
	sb.ServiceID = string(values[1:3])
	sb.LogicalTerminal = string(values[3:15])
	sb.SessionNumber = string(values[15:19])
	sb.SequenceNumber = string(values[19:25])
}

func (sb *Block2) loadValues(values string) {
	values = strings.ReplaceAll(values, " ", "")
	if values == "" {
		return
	}
	sb.Values = values
	sb.IO = string(values[0])
	sb.MessageType = string(values[1:4])

	if sb.IO == "I" {
		var ii *Block2Input = &Block2Input{}

		ii.DestinationAddress = string(values[4:16])

		if len(values) >= 17 {
			sb.MessagePriority = string(values[16:17])
		}
		if len(values) >= 18 {
			ii.DeliveryMonitoring = string(values[17:18])
		}
		if len(values) >= 21 {
			ii.ObsolescencePeriod = string(values[18:21])
		}

		sb.InputDescription = ii
	}

	if sb.IO == "O" {
		var oo *Block2Output = &Block2Output{
			InputTime:          string(values[4:8]),
			MIRDate:            string(values[8:14]),
			MIRLogicalTerminal: string(values[14:26]),
			MIRSessionNumber:   string(values[26:30]),
			MIRSequenceNumber:  string(values[30:36]),
			OutputDate:         string(values[36:42]),
			OutputTime:         string(values[42:46]),
		}

		if len(values) >= 47 {
			sb.MessagePriority = string(values[46:47])
		}

		sb.OutputDescription = oo
	}
}
