package captcha

type Composite struct {
	yes  *YesCaptcha
	cap  *CapMonster
	used string
}

func (this *Composite) CreateTask() (string, error) {
	if this.used == "YesCaptcha" {
		return this.yes.CreateTask()
	} else if this.used == "CapMonster" {
		return this.cap.CreateTask()
	}
	panic("no cap")
}

func (this *Composite) GetTaskResult(taskId string) (string, error) {
	if this.used == "YesCaptcha" {
		return this.yes.GetTaskResult(taskId)
	} else if this.used == "CapMonster" {
		return this.cap.GetTaskResult(taskId)
	}
	panic("no cap")
}
