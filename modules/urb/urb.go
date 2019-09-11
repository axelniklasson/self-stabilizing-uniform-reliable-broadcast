package urb

import (
	"log"
	"self-stabilizing-uniform-reliable-broadcast/modules/hbfd"
	"self-stabilizing-uniform-reliable-broadcast/modules/thetafd"
	"time"
)

func (module UrbModule) recordMark(k int) int {
	// TODO implement
	return -1
}

func (module UrbModule) update(m *Message, j int, s int, k int) {
	// check if buffer has record with identifier id
	id := Identifier{j, s}

	r := module.Buffer.Contains(id)
	if r != nil {
		// buffer contains record, only update recBy to include i
		r.RecBy[j] = true
		r.RecBy[k] = true
	} else {
		// create and add record to buffer
		r := BufferRecord{m, id, false, map[int]bool{j: true, k: true}, nil}
		module.Buffer.Records = append(module.Buffer.Records, r)
	}
}

func (module UrbModule) urbBroadcast(m *Message) {
	// wait(seq − min{seqMin[k]}k∈trusted < bufferUnitSize);
	module.Seq++
	module.update(m, module.ID, module.Seq, module.ID)
}

func (module UrbModule) urbDeliver(m *Message) {
	// TODO implement
}

// DoForever starts the algorithm and runs forever
func (module UrbModule) DoForever() {
	for {
		// line 13 - check if buffer contains stale information
		flush := false
		foundIdentifiers := map[Identifier]int{}
		for _, record := range module.Buffer.Records {
			// check if message is empty
			if record.Msg == nil {
				flush = true
				break
			} else {
				if _, ok := foundIdentifiers[record.Identifier]; ok {
					// same identifier used in multiple records
					flush = true
					break
				}

				foundIdentifiers[record.Identifier] = 1
			}

		}
		// line 14 - flush buffer if records either have empty messages or same identifiers
		if flush {
			module.Buffer = Buffer{}
		}

		// line 15 - remove obsolete messages
		records := []BufferRecord{}
		for _, record := range module.Buffer.Records {
			j := record.Identifier.ID
			s := record.Identifier.Seq
			if contains(module.P, j) && module.recordMark(j) <= s+module.BufferUnitSize {
				records = append(records, record)
			}
		}
		module.Buffer.Records = records

		// lines 16-19 - process records that are ready to be delivered
		for _, record := range module.Buffer.Records {
			trusted := thetafd.Trusted()
			if isSubset(trusted, record.RecBy) && !record.Delivered {
				module.urbDeliver(record.Msg)
			}

			record.Delivered = record.Delivered || isSubset(trusted, record.RecBy)
			u := hbfd.HB()
			for _, pK := range module.P {
				if _, exists := record.RecBy[pK]; !exists {
					if record.PrevHB[pK] != u[pK] {
						record.PrevHB = u
						// TODO sendMSG(m,j,s)topk;
					}
				}
			}
		}

		// lines 20-21 - send gossip messages to other processors
		for _, pK := range module.P {
			// find highest sequence number in buffer for id = pK
			s := -1
			for _, record := range module.Buffer.Records {
				if record.Identifier.ID == pK && record.Identifier.Seq > s {
					s = record.Identifier.Seq
				}
			}

			// send GOSSIP(max{s : (•, id = k, seq = s, •) ∈ buffer }, recordMark(k)) to pk
			// TODO send GOSSIP message to pK
		}

		time.Sleep(time.Second * 1)
		log.Printf("One iteration of doForever() done")
	}
}
