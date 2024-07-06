package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	clickhouse_agent "github.com/oj-lab/oj-lab-platform/modules/agent/clickhouse"
)

func TestClickhouse(t *testing.T) {
	t.Log("TestClickhouse")
	conn, err := clickhouse_agent.Connect()
	if err != nil {
		t.Error(err)
	}
	t.Log(conn)
	err = conn.Exec(context.Background(), "DROP TABLE IF EXISTS example")
	if err != nil {
		t.Error(err)
	}
	err = conn.Exec(ctx, `
    CREATE TABLE IF NOT EXISTS example (
        Col1 UInt8,
		Col2 String,
		Col3 FixedString(3),
		Col4 UUID,
		Col5 Map(String, UInt8),
		Col6 Array(String),
		Col7 Tuple(String, UInt8, Array(Map(String, String))),
		Col8 DateTime
    ) Engine = Memory
`)
	if err != nil {
		t.Error(err)
	}

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO example")
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 1000; i++ {
		err := batch.Append(
			uint8(42),
			"ClickHouse",
			"Inc",
			uuid.New(),
			map[string]uint8{"key": 1},             // Map(String, UInt8)
			[]string{"Q", "W", "E", "R", "T", "Y"}, // Array(String)
			[]interface{}{ // Tuple(String, UInt8, Array(Map(String, String)))
				"String Value", uint8(5), []map[string]string{
					{"key": "value"},
					{"key": "value"},
					{"key": "value"},
				},
			},
			time.Now(),
		)
		if err != nil {
			t.Error(err)
		}
	}
	err = batch.Send()
	if err != nil {
		t.Error(err)
	}
}
