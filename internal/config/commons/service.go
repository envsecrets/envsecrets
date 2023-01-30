package commons

import "os"

var (
	EXECUTABLE, _ = os.Executable()
	HOME_DIR, _   = os.UserHomeDir()
)

/* type Service interface {
	Save(config *Project) error
	Fetch() (*Project, error)
}

var instance Service

func SetService(svc Service) {
	if instance != nil {
		panic("service instance is already set")
	}
	instance = svc
}

func GetService() Service {
	return instance
}
*/
