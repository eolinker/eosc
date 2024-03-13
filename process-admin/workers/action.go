/*
 * Copyright (c) 2024. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package workers

import "github.com/eolinker/eosc"

type actionType int

const (
	actionCreate actionType = iota
	actionSet
	actionDelete
)

type actionContent struct {
	id     string
	action actionType
	config *eosc.WorkerConfig
}

func newActionContent(action actionType, id string, config *eosc.WorkerConfig) *actionContent {
	return &actionContent{id: id, action: action, config: config}
}
