/*
 *
 * Copyright (c) 2020 Elchin Gasimov. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 *
 * This file contents is copied from vesoft-inc/nebula-go library, make necessary changes
 */

package nebula_go_sdk

type HostAddress struct {
	Host string
	Port int
}

type timezoneInfo struct {
	offset int32
	name   []byte
}
