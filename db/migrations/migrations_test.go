package migrations

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreate_NumbersAfterLatest(t *testing.T) {
	dir := t.TempDir()
	for _, f := range []string{"0001_init.sql", "0007_products.sql", "README.md"} {
		if err := os.WriteFile(filepath.Join(dir, f), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	path, err := Create(dir, "add_stock")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if filepath.Base(path) != "0008_add_stock.sql" {
		t.Errorf("created %s, want 0008_add_stock.sql", filepath.Base(path))
	}
	b, err := os.ReadFile(path)
	if err != nil || string(b) != template {
		t.Errorf("content %q, err %v", b, err)
	}
}

func TestCreate_RejectsBadNames(t *testing.T) {
	for _, name := range []string{"", "AddStock", "add stock", "../evil", "add-stock"} {
		if _, err := Create(t.TempDir(), name); err == nil {
			t.Errorf("Create accepted %q", name)
		}
	}
}

func TestEmbeddedMigrationsParse(t *testing.T) {
	// Applying the migrations to a real database is covered by the
	// integration tests in internal/platform/db.
	entries, err := files.ReadDir(".")
	if err != nil || len(entries) == 0 {
		t.Fatalf("no embedded migrations: %v", err)
	}
	for _, e := range entries {
		if !fileName.MatchString(e.Name()) {
			t.Errorf("embedded file %s does not follow NNNN_name.sql", e.Name())
		}
	}
}
