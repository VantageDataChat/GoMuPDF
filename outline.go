package gomupdf

// Outline represents a document outline (bookmark/TOC) entry.
// Corresponds to PyMuPDF's fitz.Outline.
type Outline struct {
	Title string
	URI   string
	Page  int
	Down  *Outline // First child
	Next  *Outline // Next sibling
	IsOpen bool
}

// Flatten returns a flat list of all outline entries with their levels.
func (o *Outline) Flatten() []TOCItem {
	var items []TOCItem
	o.flatten(1, &items)
	return items
}

func (o *Outline) flatten(level int, items *[]TOCItem) {
	if o == nil {
		return
	}
	*items = append(*items, TOCItem{
		Level: level,
		Title: o.Title,
		Page:  o.Page,
	})
	if o.Down != nil {
		o.Down.flatten(level+1, items)
	}
	if o.Next != nil {
		o.Next.flatten(level, items)
	}
}
