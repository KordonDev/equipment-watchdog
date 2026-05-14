// Package audit provides structured logging for sensitive mutations
// (especially equipment <-> member assignments) so we can trace exactly
// when and why a member loses or gains a piece of equipment.
//
// All log lines share the prefix "[AUDIT]" and are emitted as a sequence
// of key=value pairs that are easy to grep and to parse.
//
// Example:
//
//	[AUDIT] event=equipment.unassign actor=mmustermann equipmentId=42 type=jacket memberIdBefore=7 memberIdAfter=0 caller=members.removeEquipmentFromMember
package audit

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/models"
)

// Field is a single key/value pair in the structured log line.
type Field struct {
	Key   string
	Value interface{}
}

// F is a tiny helper to build a Field with less ceremony.
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// Actor extracts the username of the calling user from the gin context.
// Falls back to "anonymous" if no username is set (e.g. during background
// jobs or tests).
func Actor(c *gin.Context) string {
	if c == nil {
		return "system"
	}
	if u := c.GetString("username"); u != "" {
		return u
	}
	return "anonymous"
}

// Log writes a single structured audit line for the given event.
// `actor` should be the username of the user causing the change (use Actor()).
// Additional context is supplied via Fields.
func Log(event, actor string, fields ...Field) {
	var b strings.Builder
	b.WriteString("[AUDIT] event=")
	b.WriteString(event)
	b.WriteString(" actor=")
	b.WriteString(quoteIfNeeded(actor))

	for _, f := range fields {
		b.WriteByte(' ')
		b.WriteString(f.Key)
		b.WriteByte('=')
		b.WriteString(formatValue(f.Value))
	}
	log.Info(b.String())
}

// LogCtx is a convenience wrapper that pulls the actor from the gin context.
func LogCtx(c *gin.Context, event string, fields ...Field) {
	Log(event, Actor(c), fields...)
}

// EquipmentSnapshot returns a deterministic, compact representation of a
// member's currently assigned equipment, suitable for before/after logging.
// The format is "type:id:size,type:id:size,...".
func EquipmentSnapshot(m *models.Member) string {
	if m == nil || len(m.Equipment) == 0 {
		return "none"
	}
	parts := make([]string, 0, len(m.Equipment))
	for t, e := range m.Equipment {
		if e == nil {
			parts = append(parts, fmt.Sprintf("%s:nil", t))
			continue
		}
		parts = append(parts, fmt.Sprintf("%s:%d:%s", t, e.Id, e.Size))
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

// EquipmentSnapshotFromList builds the same snapshot from a flat list.
func EquipmentSnapshotFromList(equipment []*models.Equipment) string {
	if len(equipment) == 0 {
		return "none"
	}
	parts := make([]string, 0, len(equipment))
	for _, e := range equipment {
		if e == nil {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s:%d:%s", e.Type, e.Id, e.Size))
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

// EquipmentBrief is a one-liner for a single equipment record.
func EquipmentBrief(e *models.Equipment) string {
	if e == nil {
		return "nil"
	}
	return fmt.Sprintf("id=%d type=%s memberId=%d size=%s rc=%s",
		e.Id, e.Type, e.MemberID, e.Size, e.RegistrationCode)
}

func formatValue(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return "nil"
	case string:
		return quoteIfNeeded(val)
	case fmt.Stringer:
		return quoteIfNeeded(val.String())
	default:
		return quoteIfNeeded(fmt.Sprintf("%v", val))
	}
}

func quoteIfNeeded(s string) string {
	if s == "" {
		return `""`
	}
	if strings.ContainsAny(s, " \t\"") {
		return fmt.Sprintf("%q", s)
	}
	return s
}
