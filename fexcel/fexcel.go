package fexcel

import (
	"fmt"
)

const Version = "2.0.0-beta.8"

func Logo() string {
	return fmt.Sprintf(`  __                  _
 / _|                | |
| |_ _____  _____ ___| |
|  _/ _ \ \/ / __/ _ \ |
| ||  __/>  < (_|  __/ |
|_| \___/_/\_\___\___|_|
                  v%s

by ONE Robotics Company
www.onerobotics.com

`, Version)
}
