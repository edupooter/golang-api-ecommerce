package model

type Cart struct {
	ID    string     `json:"id" example:"cart-123"`
	Items []CartItem `json:"items"`
}

func (c *Cart) AddItem(it CartItem) {
	for i := range c.Items {
		if c.Items[i].ProductID == it.ProductID {
			c.Items[i].Quantity += it.Quantity
			return
		}
	}
	c.Items = append(c.Items, it)
}

func (c *Cart) RemoveItem(productID int64) {
	for i := range c.Items {
		if c.Items[i].ProductID == productID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return
		}
	}
}
