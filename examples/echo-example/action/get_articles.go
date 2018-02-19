package action

type GetArticles struct {
}

func (ctrl GetArticles) Action() (interface{}, error) {
	return DB.All(), nil
}
