package common

// BaseCommand menyediakan implementasi dasar untuk beberapa metode Command
type BaseCommand struct {
	Name        string
	Description string
	Usage       string
	Category    string
}

// GetName mengembalikan nama command
func (c *BaseCommand) GetName() string {
	return c.Name
}

// GetDescription mengembalikan deskripsi command
func (c *BaseCommand) GetDescription() string {
	return c.Description
}

// GetCategory mengembalikan kategori command
func (c *BaseCommand) GetCategory() string {
	return c.Category
}

// GetUsage mengembalikan contoh penggunaan command
func (c *BaseCommand) GetUsage() string {
	if c.Usage != "" {
		return c.Usage
	}
	return "!" + c.Name
}
