package session

func (s *Service) routes() {
	s.Router.Methods("POST").Name("CreateSession").
		Path("/session").HandlerFunc(s.handleCreateSession())
	s.Router.Methods("PATCH").Name("UpdateSessionAttributes").
		Path("/session").HandlerFunc(s.handleUpdateSessionAttributes())
	s.Router.Methods("DELETE").Name("DropSession").
		Path("/session/{id}").HandlerFunc(s.handleDropSession())
	s.Router.Methods("GET").Name("GetSessionAttributes").
		Path("/session/{id}").HandlerFunc(s.handleGetSessionAttributes())
}
