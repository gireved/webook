package loggerxxx

import "go.uber.org/zap"

// 在 main 函数里面 Logger = xxx
var Logger *zap.Logger

var CommonLogger *zap.Logger
var SensitiveLogger *zap.Logger
