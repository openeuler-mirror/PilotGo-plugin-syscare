/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Apr 3 09:05:52 2024 +0800
 */
package service

import (
	"sync"
	"time"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/config"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/dao"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

type Task struct {
	TaskId           string   `json:"taskId"`
	IP               string   `json:"ip"`
	BuildKernel      string   `json:"buildKernel"`
	BuildDebugInfo   string   `json:"buildDebugInfo"`
	PatchDescription string   `json:"patchDescription"`
	PatchVersion     string   `json:"patchVersion"`
	PatchRelease     string   `json:"patchRelease"`
	PatchType        string   `json:"patchType"`
	Patchs           []string `json:"patchs"`
}

type Agent struct {
	IP                string
	Tasks             chan *Task
	CurrentTasksCount int // 当前正在执行的任务数量
	mutex             sync.Mutex
}

type TaskQueue struct {
	tasks   []*Task
	running bool
	agents  map[string]*Agent // 使用map来动态管理Agent，键为Agent的IP
	mutex   sync.Mutex
}

var MyTask *TaskQueue

func CreateTaskQueue() {
	MyTask = newTaskQueue()
	MyTask.start()
}

func (q *TaskQueue) Enqueue(task *Task) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.tasks = append(q.tasks, task)
}

func newTaskQueue() *TaskQueue {
	return &TaskQueue{
		running: false,
		agents:  map[string]*Agent{},
	}
}

func (q *TaskQueue) start() {
	if !q.running {
		q.running = true
		go q.dispatchTasks()
	}
}

func (q *TaskQueue) dispatchTasks() {
	for {
		if len(q.tasks) == 0 {
			time.Sleep(time.Second)
			continue
		}

		task := q.tasks[0]
		q.tasks = q.tasks[1:]

		agent, exists := q.agents[task.IP]
		if !exists {
			agent = q.createAgent(task.IP)
		}

		if agent.taskLimit() {
			agent.Tasks <- task // 如果可以执行任务，则将任务发送给 Agent
		} else {
			q.tasks = append([]*Task{task}, q.tasks...) // 如果不能执行任务，则将任务重新放回队列
		}

	}
}
func (q *TaskQueue) createAgent(ip string) *Agent {
	agent := &Agent{
		IP:    ip,
		Tasks: make(chan *Task),
	}
	q.agents[ip] = agent

	go agent.run()
	return agent
}

func (q *TaskQueue) GetTaskAgent(ip string) *Agent {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for _, agent := range q.agents {
		if agent.IP == ip {
			return agent
		}
	}
	return nil
}

func (a *Agent) run() {
	for task := range a.Tasks {
		go func(task *Task) {
			a.incrementCurrentTasksCount()
			UpdateTaskStatusToBuilding(task.TaskId)

			url := "http://" + task.IP + ":" + config.Config().AgentServer.Port + "/plugin_agent_manage/run"
			_, err := httputils.Post(url, &httputils.Params{
				Body: task,
			})
			if err != nil {
				logger.Info("下发到远程agent命令失败:%v", err.Error())
			}
		}(task)
	}
}

func (a *Agent) taskLimit() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// 如果当前正在执行的任务数量小于 maxTaskNum，则可以继续执行任务
	return a.CurrentTasksCount < dao.MaxTaskNum(a.IP)
}

func (a *Agent) incrementCurrentTasksCount() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.CurrentTasksCount++
}

func (a *Agent) DecrementCurrentTasksCount() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.CurrentTasksCount > 0 {
		a.CurrentTasksCount--
	}
}
