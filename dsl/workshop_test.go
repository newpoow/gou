package dsl

import (
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestOpenWorkshop(t *testing.T) {
	root := os.Getenv("GOU_TEST_APP_ROOT")
	workshop, err := OpenWorkshop(root)
	if err != nil {
		t.Fatal(err)
	}

	// utils.Dump(workshop)
	assert.Equal(t, 8, len(workshop.Require))
	assert.Equal(t, 16, len(workshop.Mapping))

	// utils.Dump(workshop)

	assert.Equal(t, true, workshop.Require[3].Replaced)
	assert.Equal(t, false, workshop.Require[3].Downloaded)
	assert.Equal(t, "github.com/yaoapp/demo-wms/cloud@e86eab4c8490", workshop.Require[3].URL)
	assert.Equal(t, "github.com/yaoapp/demo-wms", workshop.Require[3].Addr)
	assert.Equal(t, "github.com", workshop.Require[3].Domain)
	assert.Equal(t, "yaoapp", workshop.Require[3].Owner)
	assert.Equal(t, "demo-wms", workshop.Require[3].Repo)
	assert.Equal(t, "/cloud", workshop.Require[3].Path)
	assert.Equal(t, "demo-wms.yaoapp.cloud", workshop.Require[3].Name)
	assert.Equal(t, "demo-wms.yaoapp.cloud", workshop.Require[3].Alias)
	assert.Equal(t, "0.0.0-e86eab4c8490", workshop.Require[3].Version.String())
	assert.Equal(t, "e86eab4c8490", workshop.Require[3].Rel)

}

func TestWorkshopGetBlank(t *testing.T) {
	root := os.TempDir()
	workshop, err := OpenWorkshop(root)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(workshop.Require))
	err = workshop.Get("github.com/yaoapp/demo-wms/cloud", "wms", nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(workshop.Require))
	assert.Equal(t, 2, len(workshop.Mapping))
	assert.Equal(t, false, workshop.Require[0].Replaced)
	assert.Equal(t, false, workshop.Require[0].Downloaded)
	assert.Equal(t, "github.com/yaoapp/demo-wms/cloud@0.9.5", workshop.Require[0].URL)
	assert.Equal(t, "github.com/yaoapp/demo-wms", workshop.Require[0].Addr)
	assert.Equal(t, "github.com", workshop.Require[0].Domain)
	assert.Equal(t, "yaoapp", workshop.Require[0].Owner)
	assert.Equal(t, "demo-wms", workshop.Require[0].Repo)
	assert.Equal(t, "/cloud", workshop.Require[0].Path)
	assert.Equal(t, "demo-wms.yaoapp.cloud", workshop.Require[0].Name)
	assert.Equal(t, "wms", workshop.Require[0].Alias)
	assert.Equal(t, "0.9.5", workshop.Require[0].Version.String())
	assert.Equal(t, "0.9.5", workshop.Require[0].Rel)
}