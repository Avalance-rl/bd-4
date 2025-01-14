package standard

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


const (
    opts = "something"
)

var (
    storageOnce sync.Once
    storage *Storage
)


func TestingMain(m testing.M) {
    storageOnce.Do(func() {
        storage = InitDB(opts)
    })
    m.Run()
}


func Test_InitDBWithPool(t *testing.T) {
    t.Log("Init DB with pool")
    pgx := InitDBWithPool(opts, 10)

    
    t.Log("pgx.Stat", pgx.pool.Stat())

}


func Test_InsertProduct(t *testing.T) {
    t.Log("Succesfull storage.Insert case")

    p := Product{
        "1", 25, "Milk", "For your bones",
    }
    err := storage.InsertProduct(&p)
    assert.NoError(t, err) 
}


func Test_SelectProducts(t *testing.T) {
    t.Log("Succesfull storage.SelectProducts")
    query := `
        INSERT INTO products (price, title, description)
        VALUES ($1, $2, $3)
        RETURNING id
    `
    p := Product{
        Price: 120,
        Title: "some title",
        Description: "some description",
    }
    
	err := storage.db.QueryRow(query, p.Price, p.Title, p.Description).Scan(&p.ID)
    assert.NoError(t, err)
    
    products, err := storage.SelectProducts()
    assert.NoError(t, err)
    assert.Equal(t, products[0].Price, p.Price)
} 

func Test_GetProduct(t *testing.T) {
    t.Log("Succesfull storage.GetProduct")

    query := `
        INSERT INTO products (id, price, title, description)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
p := Product{
        ID: "999",
        Price: 120,
        Title: "some title",
        Description: "some description",
    }
    
	err := storage.db.QueryRow(query, p.ID, p.Price, p.Title, p.Description).Scan(&p.ID)
    assert.NoError(t, err)
    
    
    product, err := storage.GetProduct(p.ID)
    assert.NoError(t, err) 
    assert.Equal(t, p.Price, product.Price)

}  

func Test_UpdateProduct(t *testing.T) {
    t.Log("Succesfull storage.UpdateProduct")
    query := `
        INSERT INTO products (id, price, title, description)
        VALUES ($1, $2, $3)
        RETURNING id
    `
p:= Product{
        Price: 120,
        Title: "some title",
        Description: "some description",
    }
    
	err := storage.db.QueryRow(query, p.Price, p.Title, p.Description).Scan(&p.ID)
    assert.NoError(t, err)
    
    updateP := Product{
        Price: 50,
        Title: "Update title",
        Description: "Update description",
    }
    
    err = storage.UpdateProduct(&updateP)
    require.NoError(t, err)

    query = `
        SELECT * FROM products WHERE id = $1
    `
    var product Product
    err = storage.db.QueryRow(query, p.ID).Scan(&product)
    require.NoError(t, err)
    assert.Equal(t, updateP.ID, product.ID)        
    assert.Equal(t, updateP.Price, product.Price)
    assert.Equal(t, updateP.Description, product.Description)
    assert.Equal(t, updateP.Title, product.Title)

}

func Test_InsertProducts(t *testing.T) {
    t.Log("Succesfull storage.InsertProducts with transactions and statement")
    
    products := make([]Product, 2)
    products = append(products,
        Product{Price: 10, Title: "Some title 1", Description: "Some desc 1"},
        Product{Price: 15, Title: "Some title 2", Description: "some desc 2"},
        )

    err := storage.InsertProducts(products)
    require.NoError(t, err)
}  

