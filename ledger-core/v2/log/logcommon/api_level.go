//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package logcommon

import (
	"fmt"
	"strings"
)

type Level uint8

// NoLevel means it should be ignored
const (
	Disabled Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	NoLevel

	LogLevelCount = iota
)
const MinLevel = DebugLevel

func (l Level) IsValid() bool {
	return l > Disabled && l < NoLevel
}

func (l Level) Equal(other Level) bool {
	return l == other
}

func (l Level) String() string {
	switch l {
	case NoLevel:
		return ""
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	case Disabled:
		//
	}
	return "ignore"
}

func ParseLevel(levelStr string) (Level, error) {
	switch strings.ToLower(levelStr) {
	case Disabled.String():
		return Disabled, nil
	case DebugLevel.String():
		return DebugLevel, nil
	case InfoLevel.String():
		return InfoLevel, nil
	case WarnLevel.String():
		return WarnLevel, nil
	case ErrorLevel.String():
		return ErrorLevel, nil
	case FatalLevel.String():
		return FatalLevel, nil
	case PanicLevel.String():
		return PanicLevel, nil
	}
	return Disabled, fmt.Errorf("unknown Level String: '%s', defaulting to NoLevel", levelStr)
}