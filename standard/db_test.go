package standard

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)


func TestGetProductByID(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error creating mock db: %v", err)
    }
    defer db.Close()

    storage := &Storage{db: db}

    t.Run("successful get", func(t *testing.T) {
        rows := sqlmock.NewRows([]string{"id", "price", "title", "description"}).
            AddRow("1", "999", "Test Product", "Description")

        mock.ExpectQuery("SELECT (.+) FROM products WHERE id = \\$1").
            WithArgs("1").
            WillReturnRows(rows)

        product, err := storage.GetProduct("1")
        
        assert.NoError(t, err)
        assert.Equal(t, "1", product.ID)
        assert.Equal(t, "Test Product", product.Title)
    })

    t.Run("product not found", func(t *testing.T) {
        mock.ExpectQuery("SELECT (.+) FROM products WHERE id = \\$1").
            WithArgs("999").
            WillReturnError(sql.ErrNoRows)

        _, err := storage.GetProduct("999")

        assert.Error(t, err)
        assert.Contains(t, err.Error(), "product not found")
    })
}


