package assembler

import (
	"mq/application/dto"
)

func Return(err error) (*dto.Response, error) {
	if nil != err {
		return nil, err
	}
	return &dto.Response{Message: "ok"}, nil
}
