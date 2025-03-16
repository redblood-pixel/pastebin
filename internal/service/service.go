package service

type Service struct {
	// Set of service interfaces
	Users
	Pastes
	Permissions
}

type Users interface {
}

type Pastes interface {
}

type Permissions interface {
}
