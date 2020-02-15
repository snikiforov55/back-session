package session

func Api() string {
	return "/api/01"
}
func (s *Service) routes() {
	s.Router.Methods("POST").Name("CreateSession").
		Path(Api() + "/session").HandlerFunc(s.handleCreateSession())
	s.Router.Methods("PATCH").Name("UpdateSessionAttributes").
		Path(Api() + "/session").HandlerFunc(s.handleUpdateSessionAttributes())
	s.Router.Methods("DELETE").Name("DropSession").
		Path(Api() + "/session/{id}").HandlerFunc(s.handleDropSession())
	s.Router.Methods("GET").Name("GetSessionAttributes").
		Path(Api() + "/session/{id}").HandlerFunc(s.handleGetSessionAttributes())
}
