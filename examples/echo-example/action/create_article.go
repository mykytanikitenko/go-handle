package action

type CreateArticle struct {
	Request struct {
		Title string `json:"title" validate:"min=3,max=40,regexp=^[a-zA-Z]*$"`
		Body  string `json:"body" validate:"min=10,max=40"`
	}
}

var id = 1

func (ctrl CreateArticle) Action() (interface{}, error) {
	DB.Insert(Article{
		ID:    id,
		Title: ctrl.Request.Title,
		Body:  ctrl.Request.Body,
	})

	id++

	return nil, nil
}
