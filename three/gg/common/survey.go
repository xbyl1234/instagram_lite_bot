// Please let author have a drink, usdt trc20: TEpSxaE3kexE4e5igqmCZRMJNoDiQeWx29
// tg: @fuckins996
package common

import "github.com/AlecAivazis/survey/v2"

func SetRequire(options *survey.AskOptions) error {
	options.Validators = append(options.Validators, survey.Required)
	return nil
}
