package gcp_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/andy2kuo/AndyGameServerGo/cloud-tool/gcp"
)

func TestUpload(t *testing.T) {
	gcp.Init("sport-record-383908-0227bdade963.json")

	test_data := make(map[string]interface{})
	test_data["a"] = 1
	test_data["b"] = "2"
	test_data["c"] = 3.4
	test_data["d"] = -5.6789101112131415
	test_data["e"] = map[string]interface{}{
		"f": 1,
		"g": "2",
		"h": 3.4,
		"i": -5.6789101112131415,
	}

	bd, _ := json.MarshalIndent(test_data, "", "    ")
	t.Log(string(bd))

	err := gcp.UploadFromMemory(context.Background(), "sport-record", "test_json1.json", bd, "application/json")
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log("Upload pass")
	}

	f, _ := os.Create("test.json")
	f.Write(bd)
}
