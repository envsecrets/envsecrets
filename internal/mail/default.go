package mail

var instance Service

func SetService(svc Service) {
	if instance != nil {
		panic("service already assigned")
	}
	instance = svc
}

func GetService() Service {
	return instance
}
