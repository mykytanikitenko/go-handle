package handler

func New(pipes PipeGroup, t interface{}, converter Converter) (*handler, error) {
	if pipes == nil {
		return nil, ErrorPipesNil
	}

	if t == nil {
		return nil, ErrorTNil
	}

	if converter == nil {
		return nil, ErrorConverterNil
	}

	h := &handler{
		pipesGroup: pipes,
		t:          t,
		convertTo:  converter,
	}

	return h, h.init()
}
