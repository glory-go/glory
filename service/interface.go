package service

import "github.com/glory-go/glory/v2/config"

//go:generate mockgen -source interface.go -destination mock/interface.go

type Service interface {
	config.Component
	// Run 被调用意味着系统真正对外提供服务，该方法只会被调用一次
	Run() error
}
