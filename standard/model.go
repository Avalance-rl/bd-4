package standard

type Product struct {
	ID          string `db:"id"`
	Price       uint   `db:"price"`
	Title       string `db:"title"`
	Description string `db:"description"`
}
