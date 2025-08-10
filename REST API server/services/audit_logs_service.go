package services

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	. "github.com/go-yaaf/yaaf-common/database"
	. "github.com/go-yaaf/yaaf-common/entity"
	. "github.com/go-yaaf/yaaf-common/utils"

	"github.com/go-yaaf/yaaf-examples/rest-api/common"
	"github.com/go-yaaf/yaaf-examples/rest-api/utils"

	. "github.com/go-yaaf/yaaf-examples/rest-api/model/common"
	. "github.com/go-yaaf/yaaf-examples/rest-api/model/entities"
)

var auditLogsServiceOnce sync.Once
var auditLogsServiceInst *AuditLogsService = nil

type AuditLogsService struct {
	BaseService
	sh *common.ServiceHub // Service hub
}

// GetAuditLogsService factory function
func GetAuditLogsService(sh *common.ServiceHub) *AuditLogsService {
	auditLogsServiceOnce.Do(func() {
		if auditLogsServiceInst == nil {
			auditLogsServiceInst = &AuditLogsService{
				BaseService: BaseService{ServiceName: "AuditLogsService"},
				sh:          sh,
			}
		}
	})
	return auditLogsServiceInst
}

// Create new audit log entry in the system
func (s *AuditLogsService) Create(td *TokenData, entity Entity) (Entity, error) {

	ent := entity.(*AuditLog)

	// Override system fields,
	ent.Id = utils.TokenUtils().NanoID()
	ent.CreatedOn = Now()
	ent.UpdatedOn = Now()
	ent.Props = nil

	if updated, er := s.sh.Database.Insert(ent); er != nil {
		return nil, er
	} else {
		s.auditLog(td, ent, actionCreate, nil, updated)
		return updated, nil
	}
}

// Get single audit log entry by id
func (s *AuditLogsService) Get(td *TokenData, id string) (Entity, error) {
	entry, err := s.sh.Database.Get(NewAuditLog, id)
	if err != nil {
		return nil, fmt.Errorf("[%s]::Get: %v", s.ServiceName, err)
	}
	entry.(*AuditLog).Props = Json{}

	before := entry.(*AuditLog).BeforeChange
	after := entry.(*AuditLog).AfterChange

	bDiff, aDiff := s.extractDiff(before, after)
	entry.(*AuditLog).BeforeChange = bDiff
	entry.(*AuditLog).AfterChange = aDiff

	return entry, nil
}

// AuditLogsFindParams Query params aggregator for find commands service
type AuditLogsFindParams struct {
	From     Timestamp // From timestamp
	To       Timestamp // To timestamp
	UserId   string    // Filter by User ID
	Action   string    // Filter by action
	ItemType string    // Filter by item type
	ItemId   string    // Filter by item ID
	ItemName string    // Filter by item name
	Search   string    // Filter by free search text
	Sort     string    // Sort descriptor (field name with suffix +/- for sort order)
	Page     int       // Page number for pagination
	Size     int       // Page size: number of items per page
}

// Find list of audit log entries by filter
func (s *AuditLogsService) Find(p AuditLogsFindParams) (entities []Entity, total int64, pages int, error error) {
	cb := func(in Entity) (out Entity) {
		in.(*AuditLog).Props = Json{}
		return in
	}
	if entities, total, error = s.sh.Database.Query(NewAuditLog).
		Range("createdOn", p.From, p.To).
		MatchAny(
			F("itemType").Like(p.Search),
			F("itemId").Like(p.Search),
			F("itemName").Like(p.Search),
		).
		MatchAll(
			F("createdOn").Gte(p.From).If(p.From > 0),
			F("createdOn").Lte(p.To).If(p.To > 0),
			F("userId").Eq(p.UserId),
			F("action").Eq(p.Action),
			F("itemType").Eq(p.ItemType),
			F("itemId").Like(p.ItemId),
			F("itemName").Like(p.ItemName),
		).
		Page(p.Page).
		Limit(p.Size).
		Sort(p.Sort).
		Apply(cb).
		Find(); error == nil {
		pages = s.calcPages(total, p.Size)
	}
	return
}

// Histogram creates audit log actions count over time: TimeSeries[float64]
func (s *AuditLogsService) Histogram(p AuditLogsFindParams) (Entity, error) {

	interval := 24 * time.Hour
	out, _, err := s.sh.Database.Query(NewAuditLog).
		MatchAny(
			F("itemType").Like(p.Search),
			F("itemId").Like(p.Search),
			F("itemName").Like(p.Search),
		).
		MatchAll(
			F("createdOn").Gte(p.From).If(p.From > 0),
			F("createdOn").Lte(p.To).If(p.To > 0),
			F("userId").Eq(p.UserId),
			F("action").Eq(p.Action),
			F("itemType").Eq(p.ItemType),
			F("itemId").Eq(p.ItemId),
			F("itemName").Eq(p.ItemName),
		).
		Histogram("", COUNT, "createdOn", interval)

	if err != nil {
		return nil, s.serviceError("Histogram", err)
	}

	timeSeries := HistogramTimeSeries("audit-log", p.From, p.To, interval, out)
	return timeSeries, nil
}

// Compare and get only the different keys
func (s *AuditLogsService) extractDiff(before, after string) (bDiff string, aDiff string) {

	// initial values
	bDiff = before
	aDiff = after

	// Convert string to json
	bMap := Json{}
	if len(bDiff) > 0 {
		_ = json.Unmarshal([]byte(bDiff), &bMap)
	}

	aMap := Json{}
	if len(aDiff) > 0 {
		_ = json.Unmarshal([]byte(aDiff), &aMap)
	}

	bMap, aMap = s.diffMaps(bMap, aMap)

	if bStr, er1 := json.Marshal(bMap); er1 == nil {
		bDiff = string(bStr)
	}
	if aStr, er2 := json.Marshal(aMap); er2 == nil {
		aDiff = string(aStr)
	}
	return
}

// return map difference
func (s *AuditLogsService) diffMaps(a, b Json) (Json, Json) {
	diffA := Json{}
	diffB := Json{}

	for k, vA := range a {
		vB, exists := b[k]
		if !exists {
			diffA[k] = vA
			continue
		}

		switch vA := vA.(type) {
		case map[string]any:
			if vbMap, ok := vB.(map[string]any); ok {
				subA, subB := s.diffMaps(vA, vbMap)
				if len(subA) > 0 {
					diffA[k] = subA
				}
				if len(subB) > 0 {
					diffB[k] = subB
				}
			} else {
				diffA[k] = vA
				diffB[k] = vB
			}
		default:
			if !reflect.DeepEqual(vA, vB) {
				diffA[k] = vA
				diffB[k] = vB
			}
		}
	}

	for k, vB := range b {
		if _, exists := a[k]; !exists {
			diffB[k] = vB
		}
	}

	return diffA, diffB
}
